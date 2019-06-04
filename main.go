package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//TODO: create a function that parse URL and switch according to URl
//TODO: create a Server constructor

// User represents a single user that is registered in a social tournaments service.
type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Balance uint   `json:"balance"`
}

// Server represents server with pointer to a DB.
type Server struct {
	DB *sql.DB
}

func dbConn() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "Sql_DB_tournament1"
	dbName := "tournament_db"

	db, err = sql.Open(dbDriver, fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName))
	if err != nil {
		return nil, fmt.Errorf("can't open db: %s", err)
	}
	return db, nil
}

func main() {
	db, err := dbConn()
	if err != nil {
		log.Print(err)
		return
	}
	defer db.Close()
	s := Server{
		DB: db,
	}
	http.HandleFunc("/user", s.addUser)
	http.HandleFunc("/user/", s.getUser)
	err = http.ListenAndServe("localhost:9000", nil)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) addUser(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.NotFound(writer, req)
		return
	}
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "cannot decode json: %s", err)
		return
	}
	insert, err := s.DB.Exec("INSERT INTO user(name,balance) VALUES(?,?)", user.Name, 0)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "could not add user: %s", err)
		return
	}
	user.ID, err = insert.LastInsertId()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(writer).Encode(struct {
		ID int64 `json:"id"`
	}{
		ID: user.ID,
	})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "can't encode json: %s\n", err)
		return
	}

}

func (s *Server) getUser(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.NotFound(writer, req)
		return
	}
	idIndex := strings.LastIndex(req.URL.Path, "/") // idIndex can't be -1
	id, err := strconv.Atoi(req.URL.Path[idIndex+1:])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "incorrect id: %s", err)
		return
	}
	var user User
	err = s.DB.QueryRow("SELECT id, name,balance FROM user WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Balance)
	if err == sql.ErrNoRows {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(writer).Encode(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "can't encode json: %s\n", err)
		return
	}
}
