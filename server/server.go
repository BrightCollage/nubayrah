package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewServer() error {

	// Create multiplexer/router for the server
	r := chi.NewRouter()

	// Include CORS handler to fix CORS error in API calls.
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Use Logger for REST API request logging.
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		index, err := os.ReadFile("../static/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(index)
	})

	r.Post("/books", handleImportBook)

	r.Get("/books", handleGetAllBooks)

	r.Get("/books/{bookID}", handleGetBook)

	r.Delete("/books/{bookID}", handleDeleteBook)

	return http.ListenAndServe("localhost:5050", r)

}
