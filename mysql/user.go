package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// User represents a single user that is registered in a social tournaments service.
type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Balance uint   `json:"balance"`
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
