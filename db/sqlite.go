package db

import (
	"log"
	"nubayrah/api/book"
	"path/filepath"

	config "github.com/spf13/viper"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

func OpenDatabase() (*gorm.DB, error) {

	// Use stdlib to open a connection to postgres db.
	db_path := filepath.Join(config.GetString("db_path"))
	log.Printf("Using databse from: %v", db_path)
	DB, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})

	// Go requires DB to be used or else it complains.
	if err != nil {
		return DB, err
	}

	DB.AutoMigrate(&book.Book{})

	return DB, err
}
