package database

import (
	"gopkg.in/jackc/pgx.v2"

	"github.com/EgorAist/database_forum/internal/models"
)

const (
	selectForum = `
	SELECT user_nick, slug, title, thread_count, post_count 
	FROM forum
	WHERE slug = $1
	`
)

func SelectForum(db *pgx.ConnPool, forum *models.Forum) error {
	err := scanForum(db.QueryRow(
		selectForum,
		forum.Slug,
	), forum)

	if err == pgx.ErrNoRows {
		return ErrForumNotFound
	}

	return err
}
