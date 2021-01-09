package database

import (
	"database/sql"
	"errors"

	"github.com/EgorAist/database_forum/internal/models"
	pgx "gopkg.in/jackc/pgx.v2"
)

var (
	ErrPostConflict = errors.New("PC")
	ErrPostNotFound = errors.New("PN")
)

// id, author, created, edited, message, parent_id, thread_id, forum_slug
func scanPostRows(r *pgx.Rows, post *models.Post) error {
	parent := sql.NullInt64{}
	err := r.Scan(&post.ID, &post.Author, &post.Created, &post.IsEdited,
		&post.Message, &parent, &post.Thread, &post.Forum)

	if parent.Valid {
		post.Parent = parent.Int64
	} else {
		post.Parent = 0
	}
	return err
}

func scanPost(r *pgx.Row, post *models.Post) error {
	parent := sql.NullInt64{}
	err := r.Scan(&post.ID, &post.Author, &post.Created, &post.IsEdited,
		&post.Message, &parent, &post.Thread, &post.Forum)

	if parent.Valid {
		post.Parent = parent.Int64
	} else {
		post.Parent = 0
	}
	return err
}
