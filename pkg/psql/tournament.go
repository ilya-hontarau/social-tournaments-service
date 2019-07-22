package psql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/illfate/social-tournaments-service/pkg/sts"
	"github.com/pkg/errors"
)

// AddTournament adds tournament with passed name and deposit. Return id of this tournament.
func (db *DB) AddTournament(ctx context.Context, name string, deposit uint64) (int64, error) {
	var id int64
	err := db.conn.QueryRowContext(ctx, `
 INSERT INTO tournaments (name,deposit)
 	      VALUES ($1, $2)
  RETURNING id`, name, deposit).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "couldn't add tournament")
	}
	return id, nil
}

// GetTournament returns tournament with passed id. If tournament isn't found,
// function returns ErrNotFound.
func (db *DB) GetTournament(ctx context.Context, id int64) (*sts.Tournament, error) {
	var (
		users    sql.NullString
		winner   sql.NullInt64
		finished bool
		t        sts.Tournament
	)
	err := db.conn.QueryRowContext(ctx, `
       SELECT id, name, deposit, prize, winner, finished, json_agg(user_id)
         FROM tournaments as t
LEFT JOIN participants as p on t.id = p.tournament_id
      WHERE t.id = $1
GROUP BY t.id`, id).
		Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize, &winner, &finished, &users)
	if err == sql.ErrNoRows {
		return nil, sts.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(users.String), &t.Users)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't unmarshal json")
	}
	if finished {
		if !winner.Valid {
			return nil, errors.New("no winner")
		}
		t.Winner = winner.Int64
	}
	return &t, nil
}

// JoinTournament adds user with passed userID to tournament with passed tournamentID.
// If tournament or user isn't found, function returns ErrNotFound.
func (db *DB) JoinTournament(ctx context.Context, tournamentID, userID int64) error {
	tx, err := db.conn.BeginTxx(ctx, nil)
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
WHERE  id = $1`, tournamentID).
		Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize, &finished)
	if err == sql.ErrNoRows {
		return sts.ErrNotFound
	}
	if err != nil {
		return errors.Wrap(err, "couldn't load tournament")
	}
	if finished {
		return errors.New("tournament has finished")
	}

	update, err := tx.ExecContext(ctx, `
UPDATE users
         SET balance = balance - $1
 WHERE id = $2`, t.Deposit, userID)
	if err != nil {
		return errors.Wrap(err, "couldn't update user balance: %s")
	}
	rows, err := update.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "couldn't process user update")
	}
	if rows == 0 {
		return sts.ErrNotFound
	}

	_, err = tx.ExecContext(ctx, `
UPDATE tournaments
         SET prize = prize + deposit
 WHERE id = $1`, tournamentID)
	if err != nil {
		return errors.Wrap(err, "couldn't increase tournament prize")
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO	participants(user_id, tournament_id)
         VALUES ($1, $2)`, userID, tournamentID)
	if err != nil {
		return errors.Wrap(err, "couldn't add user to tournament")
	}
	return tx.Commit()
}
