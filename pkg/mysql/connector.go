package mysql

import (
	"database/sql"
	"errors"
	"fmt"
)

// Connector provides connection to db.
type Connector struct {
	DB *sql.DB // DB field has to be an unexported
}

// ErrNotFound is returned when item hasn't been found in db.
var ErrNotFound = errors.New("not found")

// New constructs new connection to db.
func New(dbUser, dbPass, dbName string) (*Connector, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName))
	if err != nil {
		return nil, fmt.Errorf("can't open db: %s", err)
	}
	return &Connector{
		DB: db,
	}, nil
}

// Close shuts down connection to db.
func (c *Connector) Close() {
	c.DB.Close()
}
