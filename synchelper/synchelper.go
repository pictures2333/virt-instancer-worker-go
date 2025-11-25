package synchelper

import (
	"Instancer-worker-go/config"
	"Instancer-worker-go/database"
	"Instancer-worker-go/schema"
	"Instancer-worker-go/utils"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

var once sync.Once

func Init() {
	once.Do(func() {
		var err error = nil

		if err = utils.CheckDir(config.FileDir); err != nil {
			log.Fatalf("Error while checking directory %s : %v", config.FileDir, err)
		}

		fileSyncDir := config.FileDir + "/sync"
		if err = utils.MustMkdir(fileSyncDir); err != nil {
			log.Fatalf("failed to create directory %s : %v", fileSyncDir, err)
		}
	})
}

func worker() (err error) {
	var (
		nflist []schema.FileMINIO
		oflist []database.FileLink
	)

	// GET DATA
	// get "new file list" from master
	nflist, err = getFileListFromMaster()
	if err != nil {
		return err
	}

	// get "old file list" from local db
	oflist, err = database.ReadFileLink(nil)

	// COMPARE
	// get new-created files (in master but not in db)
	new_created := find_new_create(&nflist, &oflist)
	// get old-deleted files (in db but not in master)
	old_deleted := find_old_delete(&oflist, &nflist)
	// get update list (in db and in master) (need to check last_modified)
	update := find_same(&nflist, &oflist)

	// RUN
	var wg sync.WaitGroup

	// do create
	for _, f := range *new_created {
		log.Printf("Detected new file %s/%s", f.Bucket, f.Filename)

		wg.Add(1)
		go do_create(&wg, &f)
	}
	wg.Wait()

	// do delete
	for _, f := range *old_deleted {
		log.Printf("Detected file %s deleted", f.Filename)

		wg.Add(1)
		go do_delete(&wg, &f)
	}
	wg.Wait()

	// do update
	for _, f := range *update {
		objNameReplaced := strings.ReplaceAll(f.Filename, "/", "_")
		filename := fmt.Sprintf("%s_%s", f.Bucket, objNameReplaced)

		var filelinks []database.FileLink
		if filelinks, err = database.ReadFileLink(&filename); err != nil {
			utils.Showerr(fmt.Sprintf("failed to get FileLink of file %s/%s : %v", f.Bucket, f.Filename, err), false)
			continue
		}

		if f.LastModified.After(filelinks[0].FileObj.LastModified) {
			log.Printf("Detected file %s/%s modified, updating...", f.Bucket, f.Filename)

			wg.Add(1)
			go do_update(&wg, &f)
		}
	}
	wg.Wait()

	// delete orphan FileObj (clean-up)
	do_delete_orphan_fileobj()

	return nil
}

func Worker(wg *sync.WaitGroup) {
	for {
		err := worker()
		if err != nil {
			log.Printf("Error : %v", err)
		}
		time.Sleep(time.Duration(config.SynchelperSleep) * time.Second)
	}

	//wg.Done()
}
