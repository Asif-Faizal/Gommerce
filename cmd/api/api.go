package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddress string
	db            *sql.DB
}

func NewAPIServer(listenAddress string, db *sql.DB) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		db:            db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	log.Println("Starting server on", s.listenAddress)

	return http.ListenAndServe(s.listenAddress, router)
}
