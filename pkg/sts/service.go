package sts

import (
	"context"
	"errors"
)

// Tournament represents a tournament in a social tournaments service.
type Tournament struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Deposit uint64  `json:"deposit"`
	Prize   uint64  `json:"prize"`
	Winner  int64   `json:"winner"`
	Users   []int64 `json:"users"`
}

// User represents a single user that is registered in a social tournaments service.
type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Balance uint64 `json:"balance"`
}

// ErrNotFound is returned when item hasn't been found in db.
var ErrNotFound = errors.New("not found")

type Service interface {
	//
	AddUser(ctx context.Context, name string) (int64, error)

	// If user isn't found, function has to return ErrNotFound.
	GetUser(ctx context.Context, id int64) (*User, error)

	// If user isn't found, function has to return ErrNotFound.
	DeleteUser(ctx context.Context, id int64) error

	// If user isn't found, function has to return ErrNotFound.
	AddPoints(ctx context.Context, id, points int64) error

	AddTournament(ctx context.Context, name string, deposit uint64) (int64, error)

	// If tournament isn't found, function has to return ErrNotFound.
	GetTournament(ctx context.Context, id int64) (*Tournament, error)

	// If tournament isn't found, function has to return ErrNotFound.
	JoinTournament(ctx context.Context, tournamentID, userID int64) error
}
