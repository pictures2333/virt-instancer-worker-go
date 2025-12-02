package main

import (
	"Instancer-worker-go/config"
	"Instancer-worker-go/database"
	"Instancer-worker-go/libvirtHelper"
	"Instancer-worker-go/minio"
	"Instancer-worker-go/schema"
	"Instancer-worker-go/synchelper"
	"Instancer-worker-go/vmhelper"
	"fmt"
	"log"
	"os"
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

	// init libvirtHelper
	libvirtHelper.Init()

	// init VPN
	//vpn.Init()

	wg.Add(1)
	go synchelper.Worker(&wg)

	TestArchLinux()
	//TestWindows11()

	// wait
	wg.Wait()
}

// tests
func TestArchLinux() {
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

	VMUUID := []string{
		"8b968341-c5c1-460a-8194-512786b51e57",
		//"a06c15e3-a3db-401f-bede-e30b8c87ae6c",
	}
	create := true

	for _, u := range VMUUID {
		if create {
			fake_data := schema.InstanceData{
				// master
				VMUUID: u,
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
			err = vmhelper.Delete(u)
		}
		fmt.Printf("Error : %v\n", err)
	}
}

// Windows 11
func TestWindows11() {
	files := make(map[string]schema.File)
	//files["iso"] = schema.File{
	//	Bucket:   "test",
	//	Filename: "win11.qcow2",
	//	Type:     "qcow2",
	//}
	// efi
	files["efi-loader"] = schema.File{
		Bucket:   "test",
		Filename: "OVMF_CODE.secboot.4m.fd",
		Type:     "raw",
	}
	files["efi-nvram-template"] = schema.File{
		Bucket:   "test",
		Filename: "OVMF_VARS.4m.fd",
		Type:     "raw",
	}

	vmxml, err := os.ReadFile("./challenge/challenge-win.xml")
	if err != nil {
		log.Fatal(err)
	}

	netxml, err := os.ReadFile("./challenge/network.xml")
	if err != nil {
		log.Fatal(err)
	}

	VMUUID := []string{
		"8b968341-c5c1-460a-8194-512786b51e57",
	}
	create := false

	for _, u := range VMUUID {
		if create {
			fake_data := schema.InstanceData{
				// master
				VMUUID: u,
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
			err = vmhelper.Delete(u)
		}
		fmt.Printf("Error : %v\n", err)
	}
}
