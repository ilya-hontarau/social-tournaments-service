package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//TODO: create a function that parse URL and switch according to URl

// User represents a single user that is registered in a social tournaments service.
type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Balance uint   `json:"balance"`
}

const (
	userEnvVar   = "DB_USER"
	passEnvVar   = "DB_PASS"
	dbNameEnvVar = "DB_NAME"
)

type message int

const (
	def message = iota
	take
	fund
)

// Server represents а db server in social tournaments service.
type Server struct {
	DB  *sql.DB
	msg message
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
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName))
	if err != nil {
		return nil, fmt.Errorf("can't open db: %s", err)
	}
	return &Server{DB: db}, nil
}

func main() {
	s, err := NewServer()
	if err != nil {
		log.Print(err)
		return
	}
	defer s.DB.Close()

	http.HandleFunc("/user", s.addUser)
	http.HandleFunc("/user/", s.handler)
	err = http.ListenAndServe("localhost:9000", nil)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) addUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.NotFound(w, req)
		return
	}
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "cannot decode json: %s", err)
		return
	}
	insert, err := s.DB.Exec("INSERT INTO user(name,balance) VALUES(?,?)", user.Name, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not add user: %s", err)
		return
	}
	user.ID, err = insert.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(struct {
		ID int64 `json:"id"`
	}{
		ID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "can't encode json: %s\n", err)
		return
	}

}

func (s *Server) handler(w http.ResponseWriter, req *http.Request) {
	index := strings.LastIndex(req.URL.Path, "/") // idIndex can't be -1
	word := req.URL.Path[index+1:]
	switch word {
	case "take":
		s.msg = take
	case "fund":
		s.msg = fund
	default:
		s.getUser(w, req)
		return
	}
	s.processBonus(w, req)

}

func (s *Server) processBonus(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.NotFound(w, req)
		return
	}

	id, err := strconv.Atoi((strings.Split(req.URL.Path, "/"))[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect id: %s", err)
		return
	}
	type Bonus struct {
		Points int `json:"points"`
	}
	var bonus Bonus
	err = json.NewDecoder(req.Body).Decode(&bonus)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "cannot decode json: %s", err)
		return
	}
	switch s.msg {
	case fund:
		_, err = s.DB.Exec("UPDATE user SET balance = balance + ? WHERE id = ?", bonus.Points, id)
	case take:
		_, err = s.DB.Exec("UPDATE user SET balance = balance - ? WHERE id = ?", bonus.Points, id)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not update balance: %s", err)
		return
	}
}

func (s *Server) getUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.NotFound(w, req)
		return
	}
	idIndex := strings.LastIndex(req.URL.Path, "/") // idIndex can't be -1
	id, err := strconv.Atoi(req.URL.Path[idIndex+1:])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "incorrect id: %s", err)
		return
	}
	var user User
	err = s.DB.QueryRow("SELECT id, name,balance FROM user WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Balance)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "can't encode json: %s\n", err)
		return
	}
}
