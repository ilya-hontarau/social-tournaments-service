package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/illfate/social-tournaments-service/pkg/sts"
)

// AddTournament adds tournament with passed name and deposit. Return id of this tournament.
func (c *Connector) AddTournament(ctx context.Context, name string, deposit uint64) (int64, error) {
	insert, err := c.db.ExecContext(ctx, `
 INSERT INTO tournaments (name,deposit)
 	  VALUES (?, ?)`,
		name, deposit)
	if err != nil {
		return 0, fmt.Errorf("couldn't add tournament: %s", err)
	}
	id, err := insert.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetTournament returns tournament with passed id. If tournament isn't found,
// function returns ErrNotFound.
func (c *Connector) GetTournament(ctx context.Context, id int64) (*sts.Tournament, error) {
	var (
		users    sql.NullString
		winner   sql.NullInt64
		finished bool
		t        sts.Tournament
	)
	err := c.db.QueryRowContext(ctx, `
	  SELECT id, name, deposit, prize, winner, finished, JSON_ARRAYAGG(user_id)
	    FROM tournaments
   LEFT JOIN participants ON id = tournament_id
	   WHERE id = ?
	GROUP BY id`, id).
		Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize, &winner, &finished, &users)
	if err == sql.ErrNoRows {
		return nil, sts.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(users.String), &t.Users)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal json: %s", err)
	}
	if finished {
		if !winner.Valid {
			return nil, fmt.Errorf("no winner")
		}
		t.Winner = winner.Int64
	}
	return &t, nil
}

// JoinTournament adds user with passed userID to tournament with passed tournamentID.
// If tournament or user isn't found, function returns ErrNotFound.
func (c *Connector) JoinTournament(ctx context.Context, tournamentID, userID int64) error {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var (
		finished bool
		t        sts.Tournament
	)
	err = tx.QueryRowContext(ctx, `
    SELECT id, name, deposit, prize, finished
      FROM tournaments
     WHERE id = ?`, tournamentID).
		Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize, &finished)
	if err == sql.ErrNoRows {
		return sts.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("couldn't load tournament: %s", err)
	}
	if finished {
		return errors.New("tournament has finished")
	}

	update, err := tx.ExecContext(ctx, `
    UPDATE users
       SET balance = balance - ?
     WHERE id = ?`, t.Deposit, userID)
	if err != nil {
		return fmt.Errorf("couldn't update user balance: %s", err)
	}
	rows, err := update.RowsAffected()
	if err != nil {
		return fmt.Errorf("couldn't process user update : %s", err)
	}
	if rows == 0 {
		return sts.ErrNotFound
	}

	_, err = tx.ExecContext(ctx, `
    UPDATE tournaments
       SET prize = prize + deposit
     WHERE id = ?`, tournamentID)
	if err != nil {
		return fmt.Errorf("couldn't increase tournament prize: %s", err)
	}

	_, err = tx.ExecContext(ctx, `
	INSERT INTO participants(user_id,tournament_id) 
	     VALUES (?,?)`, userID, tournamentID)
	if err != nil {
		return fmt.Errorf("couldn't add user to tournament: %s", err)
	}
	return tx.Commit()
}
