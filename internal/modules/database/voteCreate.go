package database

import (
	"github.com/EgorAist/database_forum/internal/models"
	"gopkg.in/jackc/pgx.v2"
)

const (
	insertVote = `
	INSERT INTO vote (nickname, voice, thread_id)
	VALUES ($1, $2, $3)
	ON CONFLICT ON CONSTRAINT unique_vote 
	DO UPDATE SET voice = EXCLUDED.voice;`
)

func VoteCreate(db *pgx.ConnPool, slugOrId string, t *models.Thread, v *models.Vote) error {

	if id, isID := isID(slugOrId); !isID {
		threadID, err := SelectThreadIDBySlug(db, slugOrId)
		if err != nil {
			return err
		}
		t.ID = int32(threadID)
	} else {
		t.ID = int32(id)
	}

	voteBool := voteIntToBool(v.Voice)
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("SET LOCAL synchronous_commit TO OFF")
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return txErr
		}
		return err
	}

	_, err = tx.Exec(insertVote, v.Nickname, voteBool, t.ID)
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return txErr
		}
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case pgErrForeignKeyViolation:
				return ErrThreadNotFound
			}
		}
		return err
	}

	err = threadUpdateVotesCountTx(tx, t)
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return txErr
		}
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return commitErr
	}

	return nil
}
