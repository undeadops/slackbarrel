package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/matryer/respond.v1"
)

// status endpoint
func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	//err := s.Db.VerifyCsv()
	//if err != nil {
	//	data := map[string]string{"status": "error", "msg": err.Error(), "ts": time.Now().Format(time.RFC3339)}
	//	respond.With(w, r, http.StatusFailedDependency, data)
	//		return
	//	}

	data := map[string]string{"status": "ok", "ts": time.Now().Format(time.RFC3339)}
	respond.With(w, r, http.StatusOK, data)
}

// configApp endpoint
func (s *Server) configApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appname := vars["appname"]

	config := s.Config.Apps[appname]
	respond.With(w, r, http.StatusOK, config)
}
