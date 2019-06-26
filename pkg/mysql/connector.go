package mysql

import (
	"database/sql"
	"fmt"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Connector provides connection to db.
type Connector struct {
	db *sql.DB
}

// New constructs new connection to db.
func New(dbUser, dbPass, dbName string) (*Connector, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName))
	if err != nil {
		return nil, fmt.Errorf("can't open db: %s", err)
	}
	return &Connector{
		db: db,
	}, nil
}

// Close shuts down connection to db.
func (c *Connector) Close() {
	c.db.Close()
}
