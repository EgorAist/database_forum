package database

import (
	"github.com/EgorAist/database_forum/internal/models"
	pgx "gopkg.in/jackc/pgx.v2"
)

const (
	selectUser = `
	SELECT nickname, fullname, about, email
	FROM "user"
	WHERE nickname = $1
	`

	selectUsersWithNickOrEmail = `
	SELECT nickname, fullname, about, email
	FROM "user"
	WHERE nickname = $1 OR email = $2
	`
)

func SelectUser(db *pgx.ConnPool, user *models.User) error {
	err := scanUser(db.QueryRow(selectUser, user.Nickname), user)
	if err == pgx.ErrNoRows {
		return ErrUserNotFound
	}

	return err
}

func SelectUsersWithNickOrEmail(db *pgx.ConnPool, nick, email string) (models.Users, error) {
	rows, err := db.Query(selectUsersWithNickOrEmail, nick, email)
	if err != nil {
		return nil, err
	}
	users := models.Users{}

	defer rows.Close()
	for rows.Next() {
		user := &models.User{}
		err := scanUserRows(rows, user)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

const (
	selectAllUsersByForum = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM forum_user fu
	JOIN "user" u ON fu.nickname = u.nickname
	WHERE fu.forum_slug = $1
	ORDER BY u.nickname`

	selectAllUsersByForumDesc = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM forum_user fu
	JOIN "user" u ON fu.nickname = u.nickname
	WHERE fu.forum_slug = $1
	ORDER BY u.nickname DESC`

	selectAllUsersByForumLimit = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM forum_user fu
	JOIN "user" u ON fu.nickname = u.nickname
	WHERE fu.forum_slug = $1
	ORDER BY u.nickname
	LIMIT $2`

	selectAllUsersByForumLimitDesc = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM forum_user fu
	JOIN "user" u ON fu.nickname = u.nickname
	WHERE fu.forum_slug = $1
	ORDER BY u.nickname DESC
	LIMIT $2`

	selectAllUsersByForumLimitSince = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM forum_user fu
	JOIN "user" u ON fu.nickname = u.nickname
	WHERE fu.forum_slug = $1 AND u.nickname > $2
	ORDER BY u.nickname
	LIMIT $3`

	selectAllUsersByForumLimitSinceDesc = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM forum_user fu
	JOIN "user" u ON fu.nickname = u.nickname
	WHERE fu.forum_slug = $1 AND u.nickname < $2
	ORDER BY u.nickname DESC
	LIMIT $3`

	selectAllUsersByForumSince = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM forum_user fu
	JOIN "user" u ON fu.nickname = u.nickname
	WHERE fu.forum_slug = $1 AND u.nickname > $2
	ORDER BY u.nickname`

	selectAllUsersByForumSinceDesc = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM forum_user fu
	JOIN "user" u ON fu.nickname = u.nickname
	WHERE fu.forum_slug = $1 AND u.nickname < $2
	ORDER BY u.nickname DESC`
)

func SelectAllUsersByForum(db *pgx.ConnPool, slug string, limit int, desc bool, since string,
	users *models.Users) error {

	if isExist, err := checkForumExist(db, slug); err != nil {
		return err
	} else if !isExist {
		return ErrForumNotFound
	}

	var rows *pgx.Rows
	var err error
	if desc == true {
		if since != "" && limit > 0 {
			rows, err = db.Query(selectAllUsersByForumLimitSinceDesc, slug, since, limit)
		} else if since != "" {
			rows, err = db.Query(selectAllUsersByForumSinceDesc, slug, since)
		} else if limit > 0 {
			rows, err = db.Query(selectAllUsersByForumLimitDesc, slug, limit)
		} else {
			rows, err = db.Query(selectAllUsersByForumDesc, slug)
		}
	} else {
		if since != "" && limit > 0 {
			rows, err = db.Query(selectAllUsersByForumLimitSince, slug, since, limit)
		} else if since != "" {
			rows, err = db.Query(selectAllUsersByForumSince, slug, since)
		} else if limit > 0 {
			rows, err = db.Query(selectAllUsersByForumLimit, slug, limit)
		} else {
			rows, err = db.Query(selectAllUsersByForum, slug)
		}
	}

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		user := &models.User{}
		err := scanUserRows(rows, user)
		if err != nil {
			return err
		}

		*users = append(*users, user)
	}
	return nil
}
