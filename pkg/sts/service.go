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

	// AddUser adds user with passed name to db. It returns id of this user.
	AddUser(ctx context.Context, name string) (int64, error)

	// GetUser returns user with passed id. If user isn't found, function returns ErrNotFound.
	GetUser(ctx context.Context, id int64) (*User, error)

	// DeleteUser deletes user with passed id. If user isn't found, function returns ErrNotFound.
	DeleteUser(ctx context.Context, id int64) error

	// AddPoints adds points to user with passed id. If user isn't found, function returns ErrNotFound.
	AddPoints(ctx context.Context, id, points int64) error

	// AddTournament adds tournament with passed name and deposit. Return id of this tournament.
	AddTournament(ctx context.Context, name string, deposit uint64) (int64, error)

	// GetTournament returns tournament with passed id. If tournament isn't found,
	// function returns ErrNotFound.
	GetTournament(ctx context.Context, id int64) (*Tournament, error)

	// JoinTournament adds user with passed userID to tournament with passed tournamentID.
	// If tournament or user isn't found, function returns ErrNotFound.
	JoinTournament(ctx context.Context, tournamentID, userID int64) error
}
