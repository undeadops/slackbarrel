package server

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"barrel/pkg/configserver"
)

type Server struct {
	Logger        *log.Logger
	NextRequestID func() string
	DataDir       string
	File          string
	Config        *configserver.Config
}

func SetupServer(config *configserver.Config, datadir string) *Server {
	return &Server{
		Logger:        log.New(os.Stdout, "", log.LstdFlags),
		NextRequestID: func() string { return strconv.FormatInt(time.Now().UnixNano(), 36) },
		DataDir:       datadir,
		Config:        config,
	}
}

// Router register necessary routes and returns an instance of a router.
func (s *Server) Router() *mux.Router {
	// Setup Mux Router
	r := mux.NewRouter()
	// Load up some http server middleware
	r.Use(s.tracing)
	r.Use(s.logging)

	// Serve Core Config, Variable is app name to serve config for
	r.HandleFunc("/config/{appname:[a-zA-Z-0-9]+}", s.configApp).Methods("GET")

	// This is not secure! Will revisit
	fs := http.FileServer(http.Dir(s.DataDir))
	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", fs))
	return r
}
