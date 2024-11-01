package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Define the directory where the built React files are located
	staticDir := filepath.Join("static")

	// Handle requests by serving files from the static directory
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fs)

	// Start the server on port 8080
	log.Println("Serving at http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
