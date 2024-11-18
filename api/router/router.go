package router

import (
	"net/http"
	"nubayrah/api/book"
	"os"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/gorm"
)

const FSPATH = "./static/"

func NewRouter(db *gorm.DB) chi.Router {

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

	// Router path for base endpoint. Serves built react project.
	r.Get("/*", HandleServeClient)

	BookService := book.NewBookService(db)

	// Book object routes
	r.Route("/books", BookService.RegisterRoutes)

	return r

}

func HandleServeClient(w http.ResponseWriter, r *http.Request) {
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
	http.FileServer(http.Dir(FSPATH)).ServeHTTP(w, r)
}
