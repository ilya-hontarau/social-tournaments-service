package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/illfate/social-tournaments-service/pkg/sts"
)

func (s *Server) AddUser(w http.ResponseWriter, req *http.Request) {
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

func (s *Server) GetUser(w http.ResponseWriter, req *http.Request) {
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

func (s *Server) DeleteUser(w http.ResponseWriter, req *http.Request) {
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

func (s *Server) AddPoints(w http.ResponseWriter, req *http.Request) {
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
