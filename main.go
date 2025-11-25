package main

import (
	"Instancer-worker-go/config"
	"Instancer-worker-go/database"
	"Instancer-worker-go/minio"
	"Instancer-worker-go/synchelper"
	"Instancer-worker-go/vmhelper"
	"log"
	"sync"
)

func main() {
	// Initialize

	var wg sync.WaitGroup

	// logger
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	// environment variables
	config.Init()

	// Database
	database.Init()

	// MinIO
	minio.Init()

	// init synchelper
	synchelper.Init()

	// init vmhelper
	vmhelper.Init()

	wg.Add(1)
	go synchelper.Worker(&wg)

	//Test()

	// wait
	wg.Wait()
}

/*
func Test() {
	files := make(map[string]schema.File)
	files["iso"] = schema.File{
		Bucket:   "test",
		Filename: "archlinux-2025.10.01-x86_64.iso",
		Type:     "raw",
	}

	vmxml, err := os.ReadFile("./challenge/challenge.xml")
	if err != nil {
		log.Fatal(err)
	}

	netxml, err := os.ReadFile("./challenge/network.xml")
	if err != nil {
		log.Fatal(err)
	}

	VMUUID := "8b968341-c5c1-460a-8194-512786b51e57"
	create := false

	if create {
		fake_data := schema.InstanceData{
			// master
			VMUUID: VMUUID,
			// config (from master)
			Config: schema.ChallengeConfig{
				VMXML:        string(vmxml),
				NetworkXML:   string(netxml),
				VNC:          true,
				AliveMinutes: -1,
				Files:        files,
			},
			// defined by instance
			Network: nil,
		}
		err = vmhelper.Create(&fake_data)
		if err == nil {
			if fake_data.Config.VNC {
				fmt.Println(fake_data.Network.VM.VNC.Passwd)
			}
		}
	} else {
		err = vmhelper.Delete(VMUUID)
	}
	fmt.Printf("Error : %v\n", err)
}
*/
