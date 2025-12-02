package libvirtHelper

import (
	"log"
	"sync"
)

var (
	manager LibvirtManager
	once    sync.Once
)

func Init() {
	once.Do(func() {
		var (
			err error
		)

		// connect
		if _, err = manager.GetConnection(); err != nil {
			log.Fatalf("Failed to connect libvirt : %v", err)
		}

		log.Println("Libvirt connected")
	})
}
