package database

import (
	"errors"

	"github.com/EgorAist/database_forum/internal/models"
	pgx "gopkg.in/jackc/pgx.v2"
)

var (
	ErrForumConflict = errors.New("FC")
	ErrForumNotFound = errors.New("FN")
)

func scanForum(r *pgx.Row, f *models.Forum) error {
	return r.Scan(
		&f.User,
		&f.Slug,
		&f.Title,
		&f.Threads,
		&f.Posts,
	)
}

const (
	checkForumExistQuery = `SELECT FROM forum WHERE slug = $1`
)

func checkForumExist(db *pgx.ConnPool, slug string) (bool, error) {
	err := db.QueryRow(checkForumExistQuery, slug).Scan()
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
