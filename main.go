package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/illfate/social-tournaments-service/mysql"

	_ "github.com/go-sql-driver/mysql"
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
	mysql.Connector
}

// NewServer constructs a Server, according to existing env variables.
func NewServer() (*Server, error) {
	dbUser := os.Getenv(userEnvVar)
	if dbUser == "" {
		return nil, fmt.Errorf(`no "%s" env variable`, userEnvVar)
	}
	dbPass := os.Getenv(passEnvVar)
	dbName := os.Getenv(dbNameEnvVar)
	if dbName == "" {
		return nil, fmt.Errorf(`no "%s" env variable`, dbNameEnvVar)
	}

	db, err := mysql.New(dbUser, dbPass, dbName)
	if err != nil {
		return nil, err
	}
	r := mux.NewRouter()
	s := Server{
		Connector: *db,
		Handler:   r,
	}
	r.HandleFunc("/user", s.addUser).Methods("POST")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}", s.getUser).Methods("GET")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}", s.deleteUser).Methods("DELETE")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}/{action:(?:fund|take)}", s.processBonus).Methods("POST")
	r.HandleFunc("/tournament", s.addTournament).Methods("POST")
	r.HandleFunc("/tournament/{id:[1-9]+[0-9]*}", s.getTournament).Methods("GET")
	return &s, nil
}

func main() {
	portNum := os.Getenv(port)
	if portNum == "" {
		log.Printf(`no "%s" env variable`, port)
		return
	}

	s, err := NewServer()
	if err != nil {
		log.Print(err)
		return
	}
	defer s.Close()

	err = http.ListenAndServe(":"+portNum, s)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) addUser(w http.ResponseWriter, req *http.Request) {
	var user mysql.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "cannot decode json: %s", err)
		return
	}
	user.ID, err = s.Connector.AddUser(req.Context(), user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
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
		fmt.Fprintf(w, "can't encode json: %s\n", err)
		return
	}
}

func (s *Server) processBonus(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect id: %s", err)
		return
	}
	bonus := struct {
		Points int `json:"points"`
	}{}
	err = json.NewDecoder(req.Body).Decode(&bonus)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "cannot decode json: %s", err)
		return
	}
	if vars["action"] == "take" {
		bonus.Points = -bonus.Points
	}
	err = s.Connector.UpdateUser(req.Context(), id, bonus.Points)
	if err == mysql.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func (s *Server) getUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect id: %s", err)
		return
	}
	user, err := s.Connector.GetUser(req.Context(), id) // need to return user
	if err == mysql.ErrNotFound {
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
		fmt.Fprintf(w, "can't encode json: %s\n", err)
		return
	}
}

func (s *Server) deleteUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect id: %s", err)
		return
	}
	err = s.Connector.DeleteUser(req.Context(), id)
	if err == mysql.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *Server) addTournament(w http.ResponseWriter, req *http.Request) {
	var t mysql.Tournament
	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "cannot decode json: %s", err)
		return
	}
	t.ID, err = s.Connector.AddTournament(req.Context(), t.Name, t.Deposit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
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
		fmt.Fprintf(w, "can't encode json: %s\n", err)
		return
	}
}

func (s *Server) getTournament(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect id: %s", err)
		return
	}
	t, err := s.Connector.GetTournament(req.Context(), id)
	if err == mysql.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "can't encode json: %s\n", err)
		return
	}
}
