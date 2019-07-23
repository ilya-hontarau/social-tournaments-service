package psql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	// import psql driver
	_ "github.com/lib/pq"
)

type DB struct {
	conn *sqlx.DB
}

func New(dbUser, dbPass, dbName string) (*DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s", dbUser, dbPass, dbName))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't connect to db")
	}
	return &DB{
		conn: db,
	}, nil
}

func (db *DB) Close() {
	db.conn.Close()
}
