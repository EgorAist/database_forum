package database

import (
	"github.com/EgorAist/database_forum/internal/models"
	"gopkg.in/jackc/pgx.v2"
)

const (
	updateUserFull = `
	UPDATE "user" SET fullname = $1, about = $2, email = $3
	WHERE nickname = $4
	RETURNING nickname, fullname, about, email`

	updateUserFullnameAbout = `
	UPDATE "user" SET fullname = $1, about = $2
	WHERE nickname = $3
	RETURNING nickname, fullname, about, email
	`
	updateUserFullnameEmail = `
	UPDATE "user" SET fullname = $1, email = $2
	WHERE nickname = $3
	RETURNING nickname, fullname, about, email`

	updateUserFullname = `
	UPDATE "user" SET fullname = $1
	WHERE nickname = $2
	RETURNING nickname, fullname, about, email
	`
	updateUserAboutEmail = `
	UPDATE "user" SET about = $1, email = $2
	WHERE nickname = $3
	RETURNING nickname, fullname, about, email`

	updateUserEmail = `
	UPDATE "user" SET email = $1
	WHERE nickname = $2
	RETURNING nickname, fullname, about, email`

	updateUserAbout = `
	UPDATE "user" SET about = $1
	WHERE nickname = $2
	RETURNING nickname, fullname, about, email`
)

func UpdateUser(db *pgx.ConnPool, user *models.User, us *models.UserUpdate) error {
	row := updateUser(db, user.Nickname, us)
	if row == nil {
		return SelectUser(db, user)
	}

	scanErr := scanUser(row, user)

	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return ErrUserNotFound
		}

		if pqError, ok := scanErr.(pgx.PgError); ok {
			switch pqError.Code {
			case pgErrCodeUniqueViolation:
				return ErrUserConflict
			}
		}

		return scanErr
	}

	return nil
}

func updateUser(db *pgx.ConnPool, nickname string, updateUser *models.UserUpdate) *pgx.Row {
	var row *pgx.Row
	if updateUser.About != "" && updateUser.Email != "" && updateUser.Fullname != "" {
		row = db.QueryRow(
			updateUserFull,
			updateUser.Fullname,
			updateUser.About,
			updateUser.Email,
			nickname,
		)
	} else if updateUser.About != "" && updateUser.Email != "" {
		row = db.QueryRow(
			updateUserAboutEmail,
			updateUser.About,
			updateUser.Email,
			nickname,
		)
	} else if updateUser.Email != "" && updateUser.Fullname != "" {
		row = db.QueryRow(
			updateUserFullnameEmail,
			updateUser.Fullname,
			updateUser.Email,
			nickname,
		)
	} else if updateUser.About != "" && updateUser.Fullname != "" {
		row = db.QueryRow(
			updateUserFullnameAbout,
			updateUser.Fullname,
			updateUser.About,
			nickname,
		)
	} else if updateUser.About != "" {
		row = db.QueryRow(
			updateUserAbout,
			updateUser.About,
			nickname,
		)
	} else if updateUser.Fullname != "" {
		row = db.QueryRow(
			updateUserFullname,
			updateUser.Fullname,
			nickname,
		)
	} else if updateUser.Email != "" {
		row = db.QueryRow(
			updateUserEmail,
			updateUser.Email,
			nickname,
		)
	} else if updateUser.About == "" && updateUser.Email == "" && updateUser.Fullname == "" {
		return nil
	}

	return row
}
