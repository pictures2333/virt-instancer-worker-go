package synchelper

import (
	"Instancer-worker-go/config"
	"Instancer-worker-go/database"
	"Instancer-worker-go/minio"
	"Instancer-worker-go/schema"
	"Instancer-worker-go/utils"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func do_create(wg *sync.WaitGroup, f *schema.FileMINIO) {
	defer wg.Done()

	var err error = nil

	fileSyncDir := config.FileDir + "/sync"

	// FileLink filename
	objNameReplaced := strings.ReplaceAll(f.Filename, "/", "_")
	filename := fmt.Sprintf("%s_%s", f.Bucket, objNameReplaced)
	// FileObj filename_real
	filenameReal := fmt.Sprintf("%d_%s", f.LastModified.UnixMilli(), filename)
	filenameRealFullPath := fmt.Sprintf("%s/%s", fileSyncDir, filenameReal)

	// create FIleLink and FileObj for the file
	if err = database.CreateFileLink(
		filename,
		filenameReal, f.LastModified, "CREATING",
	); err != nil {
		utils.Showerr(fmt.Sprintf("failed to create FileLink and FileObj for file %s : %v", filename, err), false)
		return
	}
	defer func(err *error) {
		if *err != nil {
			if err := database.DeleteFileLink(filename); err != nil {
				utils.Showerr(fmt.Sprintf("failed to delete FileLink %s : %v", filename, err), true)
			}
		}
	}(&err)

	// download file
	if err = minio.Download(f.Bucket, f.Filename, filenameRealFullPath); err != nil {
		utils.Showerr(fmt.Sprintf("failed to download file %s/%s from MinIO : %v", f.Bucket, f.Filename, err), false)
		return
	}
	defer func(err *error) {
		if *err != nil {
			if err := os.Remove(filenameRealFullPath); err != nil {
				utils.Showerr(fmt.Sprintf("failed to delete file %s : %v", filenameRealFullPath, err), true)
			}
		}
	}(&err)

	// update FileObj status
	if err = database.UpdateFileObjStatus(filenameReal, "READY"); err != nil {
		utils.Showerr(fmt.Sprintf("failed to update status of FileObj %s : %v", filenameReal, err), false)
		return
	}
}

func do_delete(wg *sync.WaitGroup, f *database.FileLink) {
	defer wg.Done()

	var err error = nil

	// FileLink filename
	filename := f.Filename

	// delete FileLink from DB
	if err = database.DeleteFileLink(filename); err != nil {
		utils.Showerr(fmt.Sprintf("failed to delete FileLink %s : %v", filename, err), false)
		return
	}
}

func do_update(wg *sync.WaitGroup, f *schema.FileMINIO) {
	defer wg.Done()

	var err error = nil

	fileSyncDir := config.FileDir + "/sync"

	// FileLink filename
	objNameReplaced := strings.ReplaceAll(f.Filename, "/", "_")
	filename := fmt.Sprintf("%s_%s", f.Bucket, objNameReplaced)
	// FileObj filename_real
	filenameReal := fmt.Sprintf("%d_%s", f.LastModified.UnixMilli(), filename)
	filenameRealFullPath := fmt.Sprintf("%s/%s", fileSyncDir, filenameReal)

	// download file
	if err = minio.Download(f.Bucket, f.Filename, filenameRealFullPath); err != nil {
		utils.Showerr(fmt.Sprintf("failed to download file %s/%s from MinIO : %v", f.Bucket, f.Filename, err), false)
		return
	}
	defer func(err *error) {
		if *err != nil {
			if err := os.Remove(filenameRealFullPath); err != nil {
				utils.Showerr(fmt.Sprintf("failed to delete file %s : %v", filenameRealFullPath, err), true)
			}
		}
	}(&err)

	// create FileObj and link it to FileLink
	if err = database.CreateFileObj(
		filename,
		filenameReal, f.LastModified, "READY",
	); err != nil {
		utils.Showerr(fmt.Sprintf("failed to create FileObj and update FileLink for file %s : %v", filename, err), false)
		return
	}
}

func do_delete_orphan_fileobj() {
	// Orphan FileObj
	// - (cond.1) No FileLink links to it
	// - (cond.2) No Instance links to it
	var err error = nil

	fileSyncDir := config.FileDir + "/sync"

	// get FileObjs which do not link to any FileLink (cond.1)
	var orphans []database.FileObj
	orphans, err = database.ReadFileObjOrphan()
	if err != nil {
		utils.Showerr(fmt.Sprintf("failed to get orphan FileObjs : %v", err), false)
		return
	}

	// run
	for _, f := range orphans {
		if len(f.Placeholders) == 0 { // (cond.2)
			log.Printf("Detected orphan FileObj %s, deleting...", f.FilenameReal)
			filenameRealFullPath := fmt.Sprintf("%s/%s", fileSyncDir, f.FilenameReal)

			// delete FileObj from DB
			if err = database.DeleteFileObjOrphan(f.FilenameReal); err != nil {
				utils.Showerr(fmt.Sprintf("failed to delete FileObj %s : %v", f.FilenameReal, err), false)
				continue
			}

			// delete file
			if err = os.Remove(filenameRealFullPath); err != nil {
				utils.Showerr(fmt.Sprintf("failed to remove file %s : %v", filenameRealFullPath, err), false)
				continue
			}

		}
	}
}
