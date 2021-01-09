package database

import (
	"gopkg.in/jackc/pgx.v2"

	"github.com/EgorAist/database_forum/internal/models"
)

const (
	pgErrCodeUniqueViolation  = "23505"
	pgErrForeignKeyViolation  = "23503"
	pgErrCodeNotNullViolation = "23502"
)

const clearQuery = `TRUNCATE ONLY post, vote, thread, forum_user, forum, "user"`

func Clear(db *pgx.ConnPool) error {
	_, err := db.Exec(clearQuery)
	return err
}

const statusQuery = `SELECT (SELECT COUNT(*) FROM forum), (SELECT COUNT(*) FROM thread), (SELECT count FROM post_count), (SELECT COUNT(*) FROM "user")`

func Status(db *pgx.ConnPool, s *models.Status) error {
	return db.QueryRow(statusQuery).Scan(&s.Forum, &s.Thread, &s.Post, &s.User)
}
