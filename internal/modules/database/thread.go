package database

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/EgorAist/database_forum/internal/models"
	pgx "gopkg.in/jackc/pgx.v2"
)

var (
	ErrThreadNotFoundAuthorOrForum = errors.New("TNAF")
	ErrThreadNotFound              = errors.New("TN")
	ErrThreadConflict              = errors.New("TC")
)

// Последовательность id, slug, user_nick, created, forum_slug, title, message, votes
func scanThread(r *pgx.Row, t *models.Thread) error {
	slug := sql.NullString{}

	err := r.Scan(&t.ID, &slug, &t.Author, &t.Created, &t.Forum, &t.Title, &t.Message, &t.Votes)
	if err != nil {
		return err
	}
	if slug.Valid {
		t.Slug = slug.String
	}
	return err
}

func scanThreadRows(r *pgx.Rows, t *models.Thread) error {
	slug := sql.NullString{}
	err := r.Scan(&t.ID, &slug, &t.Author, &t.Created, &t.Forum, &t.Title, &t.Message, &t.Votes)
	if err != nil {
		return err
	}
	if slug.Valid {
		t.Slug = slug.String
	}
	return err
}

func isID(slugOrID string) (int, bool) {
	if value, err := strconv.Atoi(slugOrID); err != nil {
		return -1, false
	} else {
		return value, true
	}
}

func slugToNullable(slug string) sql.NullString {
	nullable := sql.NullString{
		String: slug,
		Valid:  true,
	}
	if slug == "" {
		nullable.Valid = false
	}

	return nullable
}

const (
	checkThreadExistAndGetIDBySlug = `
	SELECT id FROM thread WHERE slug = $1
	`

	checkThreadExistAndGetIDForumSlugBySlug = `
	SELECT id, forum_slug FROM thread WHERE slug = $1
	`

	checkThreadExistAndGetForumSlugByID = `
	SELECT forum_slug FROM thread WHERE id = $1
	`

	checkThreadExistByID = `
	SELECT FROM thread WHERE id = $1
	`
)

func ifThreadExistGetID(db *pgx.ConnPool, slug string) (int, bool, error) {
	id := -1
	err := db.QueryRow(checkThreadExistAndGetIDBySlug, slug).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return id, false, nil
		}
		return id, false, err
	}
	return id, true, nil
}

func isThreadExist(db *pgx.ConnPool, id int) (bool, error) {
	err := db.QueryRow(checkThreadExistByID, id).Scan()
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ifThreadExistAndGetFodumSlugByID(db *pgx.ConnPool, id int) (string, bool, error) {
	forum := ""
	err := db.QueryRow(checkThreadExistAndGetForumSlugByID, id).Scan(&forum)
	if err != nil {
		if err == pgx.ErrNoRows {
			return forum, false, nil
		}
		return forum, false, err
	}
	return forum, true, nil
}

func ifThreadExistAndGetIDForumSlugBySlug(db *pgx.ConnPool, slug string) (string, int, bool, error) {
	id := -1
	forum := ""
	err := db.QueryRow(checkThreadExistAndGetIDForumSlugBySlug, slug).Scan(&id, &forum)
	if err != nil {
		if err == pgx.ErrNoRows {
			return forum, id, false, nil
		}
		return forum, id, false, err
	}
	return forum, id, true, nil
}
