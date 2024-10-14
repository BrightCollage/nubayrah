package db

import (
	"fmt"
	"nubayrah/api/book"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "nubayrah"
	password = "nubayrah"
	dbname   = "postgres"
)

func OpenDatabase() (*gorm.DB, error) {
	var err error
	var DB *gorm.DB
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname)

	// Use stdlib to open a connection to postgres db.
	DB, err = gorm.Open(postgres.Open(psqlInfo))

	// Go requires DB to be used or else it complains.
	if err != nil {
		return DB, err
	}

	DB.AutoMigrate(&book.Book{})
	// err = upgradeDB()

	return DB, err
}

func CloseDatabase() {
	// DB.Close()
}

// func upgradeDB() error {
// 	currentVersion, err := getVersion()
// 	if err != nil {
// 		return err
// 	}

// 	ls, err := MIGRATIONS.ReadDir("migrations")
// 	if err != nil {
// 		return err
// 	}
// 	lastestVersion := len(ls)
// 	if currentVersion >= lastestVersion {
// 		return nil
// 	}

// 	log.Printf("upgrading database %d -> %d", currentVersion, lastestVersion)

// 	for i := currentVersion + 1; i <= lastestVersion; i++ {
// 		log.Printf("applying database upgrade %d", i)
// 		cmd, err := getMigrationSQL(i)
// 		if err != nil {
// 			log.Printf("error reading database upgrade %v", err)
// 			return err
// 		}
// 		_, err = DB.Exec(cmd)
// 		if err != nil {
// 			log.Printf("error applying database upgrade %v", err)
// 			return err
// 		}
// 		err = setVersion(i)
// 		if err != nil {
// 			log.Printf("error applying database version %v", err)
// 			return err
// 		}
// 	}
// 	return nil
// }

// func getVersion() (int, error) {
// 	row := DB.QueryRow("SELECT userVersion FROM dbmetadata WHERE id = 0;")

// 	var version int
// 	err := row.Scan(&version)

// 	if err == nil {
// 		return version, nil
// 	}

// 	pqErr := new(pq.Error)
// 	if errors.As(err, &pqErr) && pqErr.Message == "relation \"dbmetadata\" does not exist" {
// 		return 0, nil
// 	}
// 	return 0, err

// }

// func setVersion(v int) error {
// 	q := fmt.Sprintf("UPDATE dbmetadata SET userVersion = %d WHERE id = 0;", v)
// 	_, err := DB.Exec(q)

// 	return err
// }

// //go:embed migrations/*
// var MIGRATIONS embed.FS

// func getMigrationSQL(version int) (string, error) {
// 	fname := fmt.Sprintf("migrations/%d.sql", version)
// 	b, err := MIGRATIONS.ReadFile(fname)
// 	if err != nil {
// 		return "", err
// 	}

// 	cmd := string(b)
// 	return cmd, nil
// }
