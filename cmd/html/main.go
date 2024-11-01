package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

// func main() {
// 	// Define the directory where the built React files are located
// 	staticDir := filepath.Join("static")

// 	// Handle requests by serving files from the static directory
// 	fs := http.FileServer(http.Dir(staticDir))
// 	http.Handle("/", fs)

// 	// Start the server on port 8080
// 	log.Println("Serving at http://localhost:8080/")
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }

const FSPATH = "./static/"

func main() {
	fs := http.FileServer(http.Dir(FSPATH))

	http.HandleFunc("/my_api", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("API CALL")) })
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
	log.Println("Serving at http://localhost:8090/")
	http.ListenAndServe(":8090", nil)
}
