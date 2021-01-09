package database

import (
	"github.com/EgorAist/database_forum/internal/models"
	"gopkg.in/jackc/pgx.v2"
)

const (
	insertThread = `
	INSERT INTO thread (slug, user_nick, created, forum_slug, title, message) 
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id, slug, user_nick, created, forum_slug, title, message, votes
	`
)

func ThreadCreate(db *pgx.ConnPool, thread *models.Thread) error {

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

	err = scanThread(tx.QueryRow(insertThread, slugToNullable(thread.Slug), thread.Author,
		thread.Created, thread.Forum,
		thread.Title, thread.Message), thread)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return txErr
		}
		if err, ok := err.(pgx.PgError); ok {
			switch err.Code {
			case pgErrCodeNotNullViolation, pgErrForeignKeyViolation:
				return ErrThreadNotFoundAuthorOrForum
			case pgErrCodeUniqueViolation:
				err := SelectThreadBySlug(db, thread)
				if err != nil {
					return err
				}
				return ErrThreadConflict
			}
		}
		return err
	}

	err = forumUpdateThreadCount(tx, thread.Forum)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return txErr
		}
		return err
	}

	err = createForumUserTx(tx, thread.Author, thread.Forum)

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
