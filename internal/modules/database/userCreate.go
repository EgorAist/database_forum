package database

import (
	"github.com/EgorAist/database_forum/internal/models"
	"gopkg.in/jackc/pgx.v2"
)

const (
	createUser = `
	INSERT INTO "user" (nickname, fullname, about, email)
	VALUES ($1, $2, $3, $4)
	RETURNING nickname, fullname, about, email
	`

	createForumUserQuery = `
	INSERT INTO forum_user (nickname, forum_slug)
	VALUES ($1, $2)
	ON CONFLICT ON CONSTRAINT unique_forum_user DO NOTHING
	`
)

func CreateUser(db *pgx.ConnPool, user *models.User) error {
	err := scanUser(db.QueryRow(
		createUser,
		user.Nickname,
		user.Fullname,
		user.About,
		user.Email,
	), user)

	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case pgErrCodeUniqueViolation:
				return ErrUserConflict
			}
		}
		return err
	}
	return nil
}

func createForumUserTx(tx *pgx.Tx, author, forum string) error {
	_, err := tx.Exec(createForumUserQuery, author, forum)
	return err
}
