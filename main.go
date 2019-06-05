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

const userEnvVar = "DB_USER"
const passEnvVar = "DB_PASS"
const dbNameEnvVar = "DB_NAME"

// Server represents а db server in social.
type Server struct {
	DB *sql.DB
}

// NewServer constructs a Server, according to existing env variables
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
	if err := db.Ping(); err != nil {
		if dbPass == "" {
			return nil, fmt.Errorf(`no "%s" env variable: %s`, passEnvVar, err)
		}
		return nil, err
	}
	return &Server{db}, nil
}

func main() {
	s, err := NewServer()
	if err != nil {
		log.Print(err)
		return
	}
	defer s.DB.Close()

	http.HandleFunc("/user", s.addUser)
	http.HandleFunc("/user/", s.getUser)
	err = http.ListenAndServe("localhost:9100", nil)
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
