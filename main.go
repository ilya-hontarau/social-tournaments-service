package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

// User represents a single user that is registered in a social tournaments service.
type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Balance uint   `json:"balance"`
}

// Tournament represents a tournament in a social tournaments service.
type Tournament struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Deposit int64   `json:"deposit"`
	Prize   int64   `json:"prize"`
	Winner  int64   `json:"winner"`
	Users   []int64 `json:"users"`
}

const (
	userEnvVar   = "DB_USER"
	passEnvVar   = "DB_PASS"
	dbNameEnvVar = "DB_NAME"
	port         = "PORT"
)

// Server represents а server in social tournaments service.
type Server struct {
	http.Handler
	DB *sql.DB
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
	r := mux.NewRouter()
	s := Server{
		DB:      db,
		Handler: r,
	}
	r.HandleFunc("/user", s.addUser).Methods("POST")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}", s.getUser).Methods("GET")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}", s.deleteUser).Methods("DELETE")
	r.HandleFunc("/user/{id:[1-9]+[0-9]*}/{category}", s.processBonus).Methods("POST")
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
	defer s.DB.Close()

	err = http.ListenAndServe(":"+portNum, s)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) addUser(w http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "cannot decode json: %s", err)
		return
	}
	insert, err := s.DB.ExecContext(req.Context(), `
 INSERT INTO users (name, balance) 
VALUES (?, ?)`,
		user.Name, 0)
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
	if req.Method != http.MethodPost {
		http.NotFound(w, req)
		return
	}
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
	var update sql.Result

	switch vars["category"] {
	case "fund":
		update, err = s.DB.ExecContext(req.Context(), `
UPDATE users 
   SET balance = balance + ? 
 WHERE id = ?`, bonus.Points, id)
	case "take":
		update, err = s.DB.ExecContext(req.Context(), `
UPDATE users 
   SET balance = balance - ? 
 WHERE id = ?`, bonus.Points, id)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not update balance: %s", err)
		return
	}
	rows, err := update.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
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
	var user User
	err = s.DB.QueryRowContext(req.Context(), `
SELECT id, name, balance 
  FROM users 
 WHERE id = ?`, id).
		Scan(&user.ID, &user.Name, &user.Balance)
	if err == sql.ErrNoRows {
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
	delete, err := s.DB.ExecContext(req.Context(), `
DELETE 
  FROM users
 WHERE id = ?`, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rows, err := delete.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (s *Server) addTournament(w http.ResponseWriter, req *http.Request) {
	var t Tournament
	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "cannot decode json: %s", err)
		return
	}
	insert, err := s.DB.ExecContext(req.Context(), `
 INSERT INTO tournaments (name,deposit)
VALUES (?,?)`,
		t.Name, t.Deposit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not add user: %s", err)
		return
	}
	t.ID, err = insert.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
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
	var (
		t        Tournament
		finished bool
		winner   sql.NullInt64
		users    string
	)
	err = s.DB.QueryRowContext(req.Context(), `
    SELECT id, name, deposit, prize, winner, finished, JSON_ARRAYAGG(user_id)
      FROM tournaments
INNER JOIN participants ON id = tournament_id 
     WHERE id = ?
  GROUP BY id`, id).
		Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize, &winner, &finished, &users)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	err = json.Unmarshal([]byte(users), &t.Users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "can't unmarshal json: %s\n", err)
		return
	}
	if finished {
		if !winner.Valid {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Winner = winner.Int64
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
