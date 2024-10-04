//go:build localdb
// +build localdb

package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
)

func openDatabase() error {
	var err error
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DB, err = sql.Open("sqlite", filepath.Join(home, "nubayrah", "db.sqlite"))
	if err != nil {
		panic(err)
	}
	return nil
}

func CloseDatabase() error {
	return DB.Close()
}
