package router

import (
	"net/http"
	"nubayrah/api/book"
	middlewares "nubayrah/api/router/middleware"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/gorm"
)

func New(db *gorm.DB) chi.Router {

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
		index, err := os.ReadFile("../../static/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(index)
	})

	BookService := book.NewBookService(db)

	// Book object routes
	r.Route("/books", func(r chi.Router) {

		r.Use(middlewares.ContentTypeJSON)

		// Book -> Create()
		r.Post("/", BookService.HandleImportBook)

		// Book -> List()
		r.Get("/", BookService.HandleGetBooks)

		// Book with object key
		r.Route("/{id}", func(r chi.Router) {

			//  Book -> Read()
			r.Get("/", BookService.HandleGetBook)

			// Book -> Delete()
			r.Delete("/", BookService.HandleDeleteBook)

			r.Route("/cover", func(r chi.Router) {

				// GetCoverImage() returns type PNG
				r.Use(middlewares.ContentTypePNG)

				//  Book -> GetCoverImage()
				r.Get("/", BookService.HandleGetBookCover)
			})

		})

	})

	return r

}
