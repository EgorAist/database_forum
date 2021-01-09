package database

import (
	"errors"

	"github.com/EgorAist/database_forum/internal/models"
	pgx "gopkg.in/jackc/pgx.v2"
)

var (
	ErrUserConflict = errors.New("UC")
	ErrUserNotFound = errors.New("UN")
)

// Последовательность Nickname Fullname About Email
func scanUser(r *pgx.Row, user *models.User) error {
	return r.Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)
}

// Последовательность Nickname Fullname About Email
func scanUserRows(r *pgx.Rows, user *models.User) error {
	return r.Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)
}
