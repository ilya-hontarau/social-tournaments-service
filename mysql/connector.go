package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

// Connector provides connection to db.
type Connector struct {
	DB *sql.DB // DB field has to be an unexported
}

// ErrNotFound is returned when item hasn't found in db.
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

// AddUser adds user with passed name to db. It returns id of this user.
func (c *Connector) AddUser(ctx context.Context, name string) (id int64, err error) {
	insert, err := c.DB.ExecContext(ctx, `
 INSERT INTO users (name) 
VALUES (?)`,
		name)
	if err != nil {
		return 0, fmt.Errorf("could not add user: %s", err)
	}
	id, err = insert.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUser returns user with passed id.
func (c *Connector) GetUser(ctx context.Context, id int) (*User, error) {
	var user User
	err := c.DB.QueryRowContext(ctx, `
SELECT id, name, balance 
  FROM users 
 WHERE id = ?`, id).
		Scan(&user.ID, &user.Name, &user.Balance)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &user, nil
}

// DeleteUser deletes user with passed id.
func (c *Connector) DeleteUser(ctx context.Context, id int) error {
	delete, err := c.DB.ExecContext(ctx, `
	DELETE 
	  FROM users
	 WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, err := delete.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// UpdateUser updates balance user with passed id.
func (c *Connector) UpdateUser(ctx context.Context, id, points int) error {
	update, err := c.DB.ExecContext(ctx, `
	UPDATE users 
	   SET balance = balance + ? 
	 WHERE id = ?`, points, id)
	if err != nil {
		return fmt.Errorf("could not update balance: %s", err)
	}
	rows, err := update.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// AddTournament adds tournament with passed name and deposit. Return id of this tournament.
func (c *Connector) AddTournament(ctx context.Context, name string, deposit int64) (int64, error) {
	insert, err := c.DB.ExecContext(ctx, `
 INSERT INTO tournaments (name,deposit)
 	  VALUES (?, ?)`,
		name, deposit)
	if err != nil {
		return 0, fmt.Errorf("could not add user: %s", err)
	}
	id, err := insert.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetTournament returns tournament with passed id.
func (c *Connector) GetTournament(ctx context.Context, id int) (*Tournament, error) {
	var (
		users    string
		winner   sql.NullInt64
		finished bool
		t        Tournament
	)
	err := c.DB.QueryRowContext(ctx, `
    SELECT id, name, deposit, prize, winner, finished, JSON_ARRAYAGG(user_id)
      FROM tournaments
INNER JOIN participants ON id = tournament_id 
     WHERE id = ?
  GROUP BY id`, id).
		Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize, &winner, &finished, &users)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(users), &t.Users)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal json: %s", err)
	}
	if finished {
		if !winner.Valid {
			return nil, err
		}
		t.Winner = winner.Int64
	}
	return &t, nil
}
