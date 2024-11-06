package main

import (
	"fmt"
	"log"
	"net/http"
	"nubayrah/api/router"
	"nubayrah/config"
	"nubayrah/db"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

const FSPATH = "./static/"

func main() {

	// API SERVER
	log.Printf("Starting Nubayrah API Server")

	err := config.Load()
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
	addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
	log.Printf("Listening at %s", addr)
	go http.ListenAndServe(addr, router.New(DB))

	// HTML SERVER
	log.Printf("Starting Nubayrah HTML Server")
	fs := http.FileServer(http.Dir(FSPATH))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If the requested file exists then return if; otherwise return index.html (fileserver default page)
		if r.URL.Path != "/" {
			fullPath := FSPATH + strings.TrimPrefix(path.Clean(r.URL.Path), "/")
			_, err := os.Stat(fullPath)
			if err != nil {
				if !os.IsNotExist(err) {
					panic(err)
				}
				// Requested file does not exist so we return the default (resolves to index.html)
				r.URL.Path = "/"
			}
		}
		fs.ServeHTTP(w, r)
	})
	addr = fmt.Sprintf("%s:%d", viper.GetString("host"), 8090)
	log.Printf("Serving at http://%s", addr)
	http.ListenAndServe(addr, nil)
}
