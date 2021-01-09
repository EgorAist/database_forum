package database

import (
	"github.com/EgorAist/database_forum/internal/models"
	"gopkg.in/jackc/pgx.v2"
)

const (
	insertForum = `
	INSERT INTO forum (user_nick, slug, title) 
	VALUES 
	((SELECT u.nickname FROM "user" u WHERE u.nickname = $1),
	$2, $3)
	RETURNING user_nick, slug, title, thread_count, post_count`
)

func CreateForum(db *pgx.ConnPool, forum *models.Forum) error {
	err := scanForum(db.QueryRow(
		insertForum,
		forum.User,
		forum.Slug,
		forum.Title,
	), forum)

	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case pgErrCodeUniqueViolation:
				return ErrForumConflict
			case pgErrCodeNotNullViolation:
				return ErrForumNotFound
			}
		}
		return err
	}

	return nil
}
