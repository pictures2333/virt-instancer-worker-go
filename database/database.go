package database

import (
	"Instancer-worker-go/config"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func Init() {
	once.Do(func() {
		var err error

		// database := sqlite.Open(config.Database)
		database := postgres.Open(config.Database)

		// connect to database
		db, err = gorm.Open(database, &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect database : %v", err)
		}

		// migrate
		err = db.AutoMigrate(
			&Instance{}, &Port{},
			&FileLink{}, &FileObj{}, &Placeholder{},
		)
		if err != nil {
			log.Fatalf("Failed to migrate tables : %v", err)
		}

		// for sqlite - enable foreign keys
		// db.Exec("PRAGMA foreign_keys = ON")

		// success
		log.Println("Database connected")
	})
}
