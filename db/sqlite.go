package db

import (
	"log"
	"nubayrah/api/book"
	"nubayrah/config"
	"path/filepath"

	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

func OpenDatabase(config config.Configuration) (*gorm.DB, error) {

	// Use stdlib to open a connection to postgres db.
	db_path := filepath.Join(config.HomeDirectory, "nubayrah.db")

	log.Printf("Using databse from: %v", db_path)
	DB, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})

	// Go requires DB to be used or else it complains.
	if err != nil {
		return DB, err
	}

	DB.AutoMigrate(&book.Book{})

	return DB, err
}
