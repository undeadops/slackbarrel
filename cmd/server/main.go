package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"barrel/pkg/configserver"
	"barrel/pkg/server"
)

var (
	port       int
	datadir    string
	configfile string
)

func main() {

	// Parse Flags
	flag.IntVar(&port, "port", 5000, "Bind Server to Port")
	flag.StringVar(&datadir, "datadir", "./data", "Directory to Load and Serve Config From")
	flag.StringVar(&configfile, "config", "./config.yaml", "Server Config")
	flag.Parse()

	// Setup Connection timeouts
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Capture Sig Terms, send to sub goroutine
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	cs := configserver.SetupConfigServer(configfile, datadir)

	// Setup Server
	s := server.SetupServer(cs, datadir)
	router := s.Router()

	// Configure HTTP Server
	s.Logger.Printf("Server is starting...")
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + strconv.Itoa(port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Startup HTTP Server
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	s.Logger.Printf("HTTP Server is ready")

	killSwitch := <-interrupt
	switch killSwitch {
	case os.Interrupt:
		s.Logger.Printf("Got SIGINT, Shutting down...")
	case syscall.SIGTERM:
		s.Logger.Printf("Got SIGTERM, Shutting down...")
	}
	srv.Shutdown(ctx)
	s.Logger.Printf("Shutdown Server")
}
