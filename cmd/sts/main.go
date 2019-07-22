package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/illfate/social-tournaments-service/pkg/psql"

	"github.com/gorilla/mux"
	"github.com/illfate/social-tournaments-service/pkg/sts"
)

const (
	port         = "PORT"
	userEnvVar   = "DB_USER"
	passEnvVar   = "DB_PASS"
	dbNameEnvVar = "DB_NAME"
)

// Server represents а server in social tournaments service.
type Server struct {
	http.Handler
	service sts.Service
}

// NewServer constructs a Server, according to existing env variables.
func NewServer(db sts.Service) *Server {
	r := mux.NewRouter()
	s := Server{
		service: db,
		Handler: r,
	}
	r.HandleFunc("/user", s.addUser).Methods("POST")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}", s.getUser).Methods("GET")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}", s.deleteUser).Methods("DELETE")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}/{action:(?:fund|take)}", s.addPoints).Methods("POST")
	r.HandleFunc("/tournament", s.addTournament).Methods("POST")
	r.HandleFunc("/tournament/{id:[1-9]+[0-9]*}", s.getTournament).Methods("GET")
	r.HandleFunc("/tournament/{id:[1-9]+[0-9]*}/join", s.joinTournament).Methods("POST")
	return &s
}

func main() {
	dbUser := os.Getenv(userEnvVar)
	if dbUser == "" {
		log.Print(fmt.Errorf(`no "%s" env variable`, userEnvVar))
		return
	}
	dbPass := os.Getenv(passEnvVar)
	dbName := os.Getenv(dbNameEnvVar)
	if dbName == "" {
		log.Print(fmt.Errorf(`no "%s" env variable`, dbNameEnvVar))
		return
	}

	db, err := psql.New(dbUser, dbPass, dbName)
	if err != nil {
		log.Print(err)
		return
	}
	defer db.Close()

	s := NewServer(db)

	portNum := os.Getenv(port)
	if portNum == "" {
		log.Printf(`no "%s" env variable`, port)
		return
	}
	err = http.ListenAndServe(":"+portNum, s)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) addUser(w http.ResponseWriter, req *http.Request) {
	var user sts.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "couldn't decode json: %s", err)
		return
	}
	user.ID, err = s.service.AddUser(req.Context(), user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't add user: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		ID int64 `json:"id"`
	}{
		ID: user.ID,
	})
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't encode json: %s\n", err)
		return
	}
}

func (s *Server) getUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect id: %s", err)
		return
	}
	user, err := s.service.GetUser(req.Context(), id)
	if err == sts.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't encode json: %s\n", err)
		return
	}
}

func (s *Server) deleteUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect id: %s", err)
		return
	}
	err = s.service.DeleteUser(req.Context(), id)
	if err == sts.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) addPoints(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect id: %s", err)
		return
	}
	bonus := struct {
		Points int64 `json:"points"`
	}{}
	err = json.NewDecoder(req.Body).Decode(&bonus)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "couldn't decode json: %s", err)
		return
	}
	if vars["action"] == "take" {
		bonus.Points = -bonus.Points
	}
	err = s.service.AddPoints(req.Context(), id, bonus.Points)
	if err == sts.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "couldn't update user: %s", err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't update user: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func (s *Server) addTournament(w http.ResponseWriter, req *http.Request) {
	var t sts.Tournament
	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "couldn't decode json: %s", err)
		return
	}
	t.ID, err = s.service.AddTournament(req.Context(), t.Name, t.Deposit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't add tournament: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		ID int64 `json:"id"`
	}{
		ID: t.ID,
	})
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't encode json: %s\n", err)
		return
	}
}

func (s *Server) getTournament(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect id: %s", err)
		return
	}
	t, err := s.service.GetTournament(req.Context(), id)
	if err == sts.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "couldn't get tournament: %s", err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't get tournament: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't encode json: %s\n", err)
		return
	}
}

func (s *Server) joinTournament(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tournamentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect id: %s", err)
		return
	}
	user := struct {
		ID int64 `json:"userId"`
	}{}
	err = json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "can't decode json: %s", err)
		return
	}
	err = s.service.JoinTournament(req.Context(), tournamentID, user.ID)
	if err == sts.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "couldn't join tournament: %s", err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't join tournament: %s", err)
		return
	}
}
