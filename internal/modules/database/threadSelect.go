package database

import (
	"github.com/EgorAist/database_forum/internal/models"
	pgx "gopkg.in/jackc/pgx.v2"
)

const (
	selectThreadByID = `
	SELECT id, slug, user_nick, created, forum_slug, title, message, votes
	FROM thread 
	WHERE id = $1
	`

	selectThreadBySlug = `
	SELECT id, slug, user_nick, created, forum_slug, title, message, votes
	FROM thread 
	WHERE slug = $1
	`

	selectThreadIDBySlug = `
	SELECT id
	FROM thread 
	WHERE slug = $1
	`
)

func SelectThreadBySlugOrID(db *pgx.ConnPool, slugOrID string, t *models.Thread) error {
	if id, isID := isID(slugOrID); !isID {
		t.Slug = slugOrID
		return SelectThreadBySlug(db, t)
	} else {
		t.ID = int32(id)
		return SelectThreadByID(db, t)
	}
}

func SelectThreadByID(db *pgx.ConnPool, t *models.Thread) error {
	err := scanThread(db.QueryRow(selectThreadByID, t.ID), t)

	if err == pgx.ErrNoRows {
		return ErrThreadNotFound
	}

	return err
}

func SelectThreadBySlug(db *pgx.ConnPool, t *models.Thread) error {
	err := scanThread(db.QueryRow(selectThreadBySlug, t.Slug), t)

	if err == pgx.ErrNoRows {
		return ErrThreadNotFound
	}

	return err
}

func SelectThreadIDBySlug(db *pgx.ConnPool, slug string) (int, error) {
	id := -1
	err := db.QueryRow(selectThreadIDBySlug, slug).Scan(&id)
	if err == pgx.ErrNoRows {
		return 0, ErrThreadNotFound
	}
	return id, err
}

const (
	selectAllThreads = `
	SELECT id, slug, user_nick, created, forum_slug, title, message, votes
	FROM thread
	WHERE forum_slug = $1
	ORDER BY created
	`

	selectAllThreadsDesc = `
	SELECT id, slug, user_nick, created, forum_slug, title, message, votes
	FROM thread
	WHERE forum_slug = $1
	ORDER BY created DESC
	`

	selectAllThreadsLimit = `
	SELECT id, slug, user_nick, created, forum_slug, title, message, votes
	FROM thread
	WHERE forum_slug = $1
	ORDER BY created
	LIMIT $2
	`

	selectAllThreadsLimitDesc = `
	SELECT id, slug, user_nick, created, forum_slug, title, message, votes
	FROM thread
	WHERE forum_slug = $1
	ORDER BY created DESC
	LIMIT $2
	`

	selectAllThreadsSince = `
	SELECT id, slug, user_nick, created, forum_slug, title, message, votes
	FROM thread
	WHERE forum_slug = $1 AND created >= $2
	ORDER BY created
	`

	selectAllThreadsSinceDesc = `
	SELECT id, slug, user_nick, created, forum_slug, title, message, votes
	FROM thread
	WHERE forum_slug = $1 AND created <= $2
	ORDER BY created DESC
	`

	selectAllThreadsSinceLimit = `
	SELECT id, slug, user_nick, created, forum_slug, title, message, votes
	FROM thread
	WHERE forum_slug = $1 AND created >= $2
	ORDER BY created
	LIMIT $3
	`

	selectAllThreadsSinceLimitDesc = `
	SELECT id, slug, user_nick, created, forum_slug, title, message, votes
	FROM thread
	WHERE forum_slug = $1 AND created <= $2
	ORDER BY created DESC
	LIMIT $3
	`
)

func SelectAllThreadsByForum(db *pgx.ConnPool, slug string, limit int, desc bool,
	since string, ts *models.Threads) error {

	if isExist, err := checkForumExist(db, slug); err != nil {
		return err
	} else if !isExist {
		return ErrForumNotFound
	}

	var rows *pgx.Rows
	var err error
	if desc == true {
		if limit > 0 && since != "" {
			rows, err = db.Query(selectAllThreadsSinceLimitDesc, slug, since, limit)
		} else if limit > 0 {
			rows, err = db.Query(selectAllThreadsLimitDesc, slug, limit)
		} else if since != "" {
			rows, err = db.Query(selectAllThreadsSinceDesc, slug, since)
		} else {
			rows, err = db.Query(selectAllThreadsDesc, slug)
		}
	} else {
		if limit > 0 && since != "" {
			rows, err = db.Query(selectAllThreadsSinceLimit, slug, since, limit)
		} else if limit > 0 {
			rows, err = db.Query(selectAllThreadsLimit, slug, limit)
		} else if since != "" {
			rows, err = db.Query(selectAllThreadsSince, slug, since)
		} else {
			rows, err = db.Query(selectAllThreads, slug)
		}
	}

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		thread := &models.Thread{}
		err := scanThreadRows(rows, thread)
		if err != nil {
			return err
		}

		*ts = append(*ts, thread)
	}

	return nil
}
