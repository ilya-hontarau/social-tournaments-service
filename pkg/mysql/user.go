package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/illfate/social-tournaments-service/pkg/sts"
)

// AddUser adds user with passed name to db. It returns id of this user.
func (c *Connector) AddUser(ctx context.Context, name string) (int64, error) {
	insert, err := c.db.ExecContext(ctx, `
 INSERT INTO users (name) 
VALUES (?)`,
		name)
	if err != nil {
		return 0, fmt.Errorf("couldn't add user: %s", err)
	}
	id, err := insert.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUser returns user with passed id. If user isn't found, function returns ErrNotFound.
func (c *Connector) GetUser(ctx context.Context, id int64) (*sts.User, error) {
	var user sts.User
	err := c.db.GetContext(ctx, &user, `
SELECT id, name, balance 
  FROM users 
 WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, sts.ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &user, nil
}

// DeleteUser deletes user with passed id. If user isn't found, function returns ErrNotFound.
func (c *Connector) DeleteUser(ctx context.Context, id int64) error {
	delete, err := c.db.ExecContext(ctx, `
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
		return sts.ErrNotFound
	}
	return nil
}

// AddPoints adds points to user with passed id. If user isn't found, function returns ErrNotFound.
func (c *Connector) AddPoints(ctx context.Context, id, points int64) error {
	update, err := c.db.ExecContext(ctx, `
	UPDATE users 
	   SET balance = balance + ? 
	 WHERE id = ?`, points, id)
	if err != nil {
		return fmt.Errorf("couldn't update balance: %s", err)
	}
	rows, err := update.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sts.ErrNotFound
	}
	return nil
}
