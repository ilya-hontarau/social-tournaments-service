package psql

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/illfate/social-tournaments-service/pkg/sts"
)

// AddUser adds user with passed name to db. It returns id of this user.
func (db *DB) AddUser(ctx context.Context, name string) (int64, error) {
	var id int64
	err := db.conn.QueryRowContext(ctx, `
 INSERT INTO users (name) 
           VALUES ($1)
   RETURNING id`, name).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "couldn't add user")
	}
	return id, nil
}

// GetUser returns user with passed id. If user isn't found, function returns ErrNotFound.
func (db *DB) GetUser(ctx context.Context, id int64) (*sts.User, error) {
	var user sts.User
	err := db.conn.GetContext(ctx, &user, `
SELECT id, name, balance 
  FROM users 
WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, sts.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get user")
	}
	return &user, nil
}

// DeleteUser deletes user with passed id. If user isn't found, function returns ErrNotFound.
func (db *DB) DeleteUser(ctx context.Context, id int64) error {
	delete, err := db.conn.ExecContext(ctx, `
	DELETE
	   FROM users
	 WHERE id = $1`, id)
	if err != nil {
		return errors.Wrap(err, "couldn't delete user")
	}
	rows, err := delete.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "couldn't get affected rows")
	}
	if rows == 0 {
		return sts.ErrNotFound
	}
	return nil
}

// AddPoints adds points to user with passed id. If user isn't found, function returns ErrNotFound.
func (db *DB) AddPoints(ctx context.Context, id, points int64) error {
	update, err := db.conn.ExecContext(ctx, `
 UPDATE users 
	      SET balance = balance + $1 
  WHERE id = $2`, points, id)
	if err != nil {
		return errors.Wrap(err, "couldn't update balance")
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
