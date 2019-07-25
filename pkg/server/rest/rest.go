package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/illfate/social-tournaments-service/pkg/sts"
)

type Server struct {
	http.Handler
	service sts.Service
}

// NewServer constructs a Server, according to existing env variables.
func New(db sts.Service) *Server {
	r := mux.NewRouter()
	s := Server{
		service: db,
		Handler: r,
	}
	r.HandleFunc("/user", s.AddUser).Methods("POST")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}", s.GetUser).Methods("GET")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}", s.DeleteUser).Methods("DELETE")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}/{action:(?:fund|take)}", s.AddPoints).Methods("POST")
	r.HandleFunc("/tournament", s.AddTournament).Methods("POST")
	r.HandleFunc("/tournament/{id:[1-9]+[0-9]*}", s.GetTournament).Methods("GET")
	r.HandleFunc("/tournament/{id:[1-9]+[0-9]*}/join", s.JoinTournament).Methods("POST")
	return &s
}
