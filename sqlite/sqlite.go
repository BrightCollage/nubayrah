package sqlite

import (
	"log"
	"nubayrah/api/book"
	"path/filepath"

	config "github.com/spf13/viper"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

func OpenDatabase(path string) (*gorm.DB, error) {

	// Open DB command
	DB, err := gorm.Open(sqlite.Open(path), &gorm.Config{})

	if err != nil {
		return DB, err
	}

	// Run Automigration
	DB.AutoMigrate(&book.Book{})

	return DB, err
}

func NewDB() *gorm.DB {
	// Use stdlib to open a connection to postgres db.
	db_path := filepath.Join(config.GetString("db_path"))
	log.Printf("Using databse from: %v", db_path)
	DB, err := OpenDatabase(db_path)
	if err != nil {
		log.Printf("error when connecting to database %v", err)
		panic(err)
	}
	return DB
}
