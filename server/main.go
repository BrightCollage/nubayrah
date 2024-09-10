package main

import (
	"log"
)

func main() {
	log.Printf("Starting Nubayrah Server")

	// Start with connecting to the Database
	err := OpenDatabase()
	if err != nil {
		log.Printf("error when connecting to postgresql database %v", err)
	}
	defer CloseDatabase()

	// Started HTTPServer
	err = NewServer()
	if err != nil {
		log.Printf("Error when trying to instantiate server %v", err)
	}

}
