package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nubayrah/api/router"
	"nubayrah/config"
	"nubayrah/sqlite"
	"os"
	"os/signal"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const FSPATH = "./static/"

func main() {

	// Setting up a signal handler to receive kill signal.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)   // Single buffer that takes os.Signal items
	signal.Notify(c, os.Interrupt) // os.Interrupt is CTRL + C signal, will notify the c channel

	// Start a goroutine to read channel and instantiate cancel functions on os trigger.
	go func() { <-c; cancel() }()

	// Execute program.
	log.Printf("Starting Nubayrah Application")
	//
	// Creates a new Main object
	m := NewMain()
	//
	// Starts the API server
	m.StartServer()
	//
	// Line to wait for CTRL-C
	<-ctx.Done()

	// Start shutting down
	log.Printf("Shutting down Server...")
	if err := m.server.Shutdown(ctx); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	log.Printf("main: done. exiting")
}

type Main struct {
	// Object to maintain main functions
	server *http.Server
	db     *gorm.DB
}

func NewMain() *Main {

	// Load configurations first.
	err := config.Load()
	if err != nil {
		panic(err)
	}

	// Get address of server to attach
	addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
	return &Main{
		server: &http.Server{
			Addr: addr,
		},
		db: sqlite.NewDB(),
	}
}

func (m *Main) StartServer() {

	// API Server to host endpoints.
	log.Printf("Starting Nubayrah API Server")

	// Attach the router to the http.Server
	m.server.Handler = router.NewRouter(m.db)
	// Start with connecting to the Database
	go func() {

		log.Printf("Listening at: http://%v", m.server.Addr)
		// always returns error. ErrServerClosed on graceful close
		if err := m.server.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

}
