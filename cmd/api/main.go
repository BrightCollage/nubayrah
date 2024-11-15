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
)

const FSPATH = "./static/"

func main() {

	// Setting up a signal handler to receive kill signal.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)   // Single buffer that takes os.Signal items
	signal.Notify(c, os.Interrupt) // os.Interrupt is CTRL + C signal, will notify the c channel

	// Start a goroutine to read channel and instantiate cancel functions on os trigger.
	go func() { <-c; cancel() }()

	log.Printf("Starting Nubayrah Application")
	// Execute program.
	//
	// Starts the API server
	Server := StartServer()
	// Starts the Client Webpage server

	// Line to wait for CTRL-C
	<-ctx.Done()

	// Start shutting down
	log.Printf("Shutting down Server...")
	if err := Server.Shutdown(ctx); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	log.Printf("main: done. exiting")
}

func StartServer() *http.Server {
	// API Server to host endpoints.

	log.Printf("Starting Nubayrah API Server")

	err := config.Load()
	if err != nil {
		panic(err)
	}

	// Start with connecting to the Database
	DB, err := sqlite.OpenDatabase()
	if err != nil {
		log.Printf("error when connecting to database %v", err)
		panic(err)
	}

	// Start API Server
	addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
	srv := &http.Server{Addr: addr, Handler: router.New(DB)}
	// http.ListenAndServe(addr, router.New(DB))
	log.Printf("Listening at %s", addr)
	go func() {
		// defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv

}
