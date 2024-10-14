package main

import (
	"fmt"
	"log"
	"net/http"
	"nubayrah/api/book"
	"nubayrah/config"
	"nubayrah/db"

	"github.com/go-chi/chi/v5"
)

func main() {
	log.Printf("Starting Nubayrah Server")

	userConfig := config.New()

	// Start with connecting to the Database
	DB, err := db.OpenDatabase()
	if err != nil {
		log.Printf("error when connecting to database %v", err)
		panic(err)
	}
	defer db.CloseDatabase()

	// Started HTTPServer
	err = NewServer(book.NewRouter(DB), userConfig)
	if err != nil {
		log.Printf("Error when trying to instantiate server %v", err)
	}

}

func NewServer(r chi.Router, config config.Configuration) error {

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Printf("Listening at %s", addr)
	return http.ListenAndServe(addr, r)

}
