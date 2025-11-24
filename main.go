package main

import (
	"Instancer-worker-go/config"
	"Instancer-worker-go/database"
	"Instancer-worker-go/minio"
	"Instancer-worker-go/synchelper"
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

	wg.Add(1)
	go synchelper.Worker(&wg)

	// wait
	wg.Wait()
}
