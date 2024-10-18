package main

import (
	"fmt"
	"log"
	"net/http"
	"nubayrah/api/router"
	"nubayrah/config"
	"nubayrah/db"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func main() {
	log.Printf("Starting Nubayrah Server")

	err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	// Start with connecting to the Database
	DB, err := db.OpenDatabase()
	if err != nil {
		log.Printf("error when connecting to database %v", err)
		panic(err)
	}

	// Started HTTPServer
	err = NewServer(router.New(DB))
	if err != nil {
		log.Printf("Error when trying to instantiate server %v", err)
	}

}

func NewServer(r chi.Router) error {

	addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
	log.Printf("Listening at %s", addr)
	return http.ListenAndServe(addr, r)

}
