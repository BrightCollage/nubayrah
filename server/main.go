package main

import (
	"log"
	"main/db"
)

func main() {
	log.Printf("Starting Nubayrah Server")

	err := OpenConfig()
	if err != nil {
		log.Printf("error opening config file %v", err)
		panic(err)
	}

	// Start with connecting to the Database
	err = db.OpenDatabase()
	if err != nil {
		log.Printf("error when connecting to database %v", err)
		panic(err)
	}
	defer db.CloseDatabase()

	// Started HTTPServer
	err = NewServer()
	if err != nil {
		log.Printf("Error when trying to instantiate server %v", err)
	}

}
