package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	host     = "db"
	port     = 5432
	user     = "nubayrah"
	password = "nubayrah"
	dbname   = "postgres"
)

func OpenDatabase() error {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname)

	// Use stdlib to open a connection to postgres db.
	DB, err = sql.Open("postgres", psqlInfo)

	// Go requires DB to be used or else it complains.
	_ = DB
	if err != nil {
		return err
	}

	return nil
}

func CloseDatabase() error {
	return DB.Close()
}
