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

func New(dbUser, dbHost, dbPass, dbName string) (*DB, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=disable", dbUser, dbPass, dbHost, dbName))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't connect to db")
	}
	return &DB{
		conn: db,
	}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
