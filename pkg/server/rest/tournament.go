package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/illfate/social-tournaments-service/pkg/sts"
)

func (s *Server) AddTournament(w http.ResponseWriter, req *http.Request) {
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

func (s *Server) GetTournament(w http.ResponseWriter, req *http.Request) {
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

func (s *Server) JoinTournament(w http.ResponseWriter, req *http.Request) {
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
