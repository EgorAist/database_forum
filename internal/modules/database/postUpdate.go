package database

import (
	"github.com/EgorAist/database_forum/internal/models"
	pgx "gopkg.in/jackc/pgx.v2"
)

const (
	updatePostFull = `
	UPDATE post SET message = $1
	WHERE id = $2
	RETURNING id, author, created, edited, message, parent_id, thread_id, forum_slug`
)

func UpdatePost(db *pgx.ConnPool, post *models.Post, pu *models.PostUpdate) error {
	var err error
	if pu.Message == "" {
		err = selectPost(db, post)
	} else {
		err = scanPost(db.QueryRow(updatePostFull, pu.Message, post.ID), post)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrPostNotFound
		}
		return err
	}
	return nil
}
