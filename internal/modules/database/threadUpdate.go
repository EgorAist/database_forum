package database

import (
	"github.com/EgorAist/database_forum/internal/models"
	pgx "gopkg.in/jackc/pgx.v2"
)

const (
	threadUpdateVotesCountQuery = `UPDATE thread t SET votes = (
	SELECT SUM(case when v.voice = true then 1 else -1 end)
	FROM vote v 
	WHERE v.thread_id=$1) WHERE t.id=$2
	RETURNING id, slug, user_nick, created, forum_slug, title, message, votes`

	threadUpdateFullByID = `
	UPDATE thread SET message = $1, title = $2
	WHERE id = $3
	RETURNING id, slug, user_nick, created, forum_slug, title, message, votes
	`

	threadUpdateMessageByID = `
	UPDATE thread SET message = $1
	WHERE id = $2
	RETURNING id, slug, user_nick, created, forum_slug, title, message, votes
	`
	threadUpdateTitleByID = `
	UPDATE thread SET title = $1
	WHERE id = $2
	RETURNING id, slug, user_nick, created, forum_slug, title, message, votes
	`

	threadUpdateFullBySlug = `
	UPDATE thread SET message = $1, title = $2
	WHERE slug = $3
	RETURNING id, slug, user_nick, created, forum_slug, title, message, votes
	`

	threadUpdateMessageBySlug = `
	UPDATE thread SET message = $1
	WHERE slug = $2
	RETURNING id, slug, user_nick, created, forum_slug, title, message, votes
	`
	threadUpdateTitleBySlug = `
	UPDATE thread SET title = $1
	WHERE slug = $2
	RETURNING id, slug, user_nick, created, forum_slug, title, message, votes
	`
)

func UpdateThread(db *pgx.ConnPool, tu *models.ThreadUpdate, slugOrID string, t *models.Thread) error {
	var row *pgx.Row
	if id, isID := isID(slugOrID); !isID {
		t.Slug = slugOrID
		row = updateThreadBySlug(db, tu, t.Slug)
		if row == nil {
			return SelectThreadBySlug(db, t)
		}
	} else {
		t.ID = int32(id)
		row = updateThreadByID(db, tu, t.ID)
		if row == nil {
			return SelectThreadByID(db, t)
		}
	}

	err := scanThread(row, t)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrThreadNotFound
		}
		return err
	}

	return nil
}

func updateThreadBySlug(db *pgx.ConnPool, tu *models.ThreadUpdate, slug string) *pgx.Row {
	var row *pgx.Row
	if tu.Message != "" && tu.Title != "" {
		row = db.QueryRow(
			threadUpdateFullBySlug,
			tu.Message,
			tu.Title,
			slug,
		)
	} else if tu.Message != "" {
		row = db.QueryRow(
			threadUpdateMessageBySlug,
			tu.Message,
			slug,
		)
	} else if tu.Title != "" {
		row = db.QueryRow(
			threadUpdateTitleBySlug,
			tu.Title,
			slug,
		)
	} else if tu.Title == "" && tu.Message == "" {
		return nil
	}

	return row
}

func updateThreadByID(db *pgx.ConnPool, tu *models.ThreadUpdate, id int32) *pgx.Row {
	var row *pgx.Row
	if tu.Message != "" && tu.Title != "" {
		row = db.QueryRow(
			threadUpdateFullByID,
			tu.Message,
			tu.Title,
			id,
		)
	} else if tu.Message != "" {
		row = db.QueryRow(
			threadUpdateMessageByID,
			tu.Message,
			id,
		)
	} else if tu.Title != "" {
		row = db.QueryRow(
			threadUpdateTitleByID,
			tu.Title,
			id,
		)
	} else if tu.Title == "" && tu.Message == "" {
		return nil
	}

	return row
}

func threadUpdateVotesCount(db *pgx.ConnPool, t *models.Thread) error {
	return scanThread(db.QueryRow(threadUpdateVotesCountQuery, t.ID, t.ID), t)
}

func threadUpdateVotesCountTx(tx *pgx.Tx, t *models.Thread) error {
	return scanThread(tx.QueryRow(threadUpdateVotesCountQuery, t.ID, t.ID), t)
}
