package db

import (
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	"log"
)

var DB *sql.DB

func OpenDatabase() error {
	err := openDatabase()
	if err != nil {
		return err
	}

	err = upgradeDB()

	return err
}

func upgradeDB() error {
	currentVersion, err := getVersion()
	if err != nil {
		return err
	}

	ls, err := MIGRATIONS.ReadDir("migrations")
	if err != nil {
		return err
	}
	lastestVersion := len(ls)
	if currentVersion >= lastestVersion {
		return nil
	}

	log.Printf("upgrading database %d -> %d", currentVersion, lastestVersion)

	for i := currentVersion + 1; i <= lastestVersion; i++ {
		log.Printf("applying database upgrade %d", i)
		cmd, err := getMigrationSQL(i)
		if err != nil {
			log.Printf("error reading database upgrade %v", err)
			return err
		}
		_, err = DB.Exec(cmd)
		if err != nil {
			log.Printf("error applying database upgrade %v", err)
			return err
		}
		err = setVersion(i)
		if err != nil {
			log.Printf("error applying database version %v", err)
			return err
		}
	}
	return nil
}

func getVersion() (int, error) {
	row := DB.QueryRow("PRAGMA user_version")

	v := 0
	err := row.Scan(&v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func setVersion(v int) error {
	q := fmt.Sprintf("PRAGMA user_version = %d", v)
	_, err := DB.Exec(q)

	return err
}

//go:embed migrations/*
var MIGRATIONS embed.FS

func getMigrationSQL(version int) (string, error) {
	fname := fmt.Sprintf("migrations/%d.sql", version)
	b, err := MIGRATIONS.ReadFile(fname)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
