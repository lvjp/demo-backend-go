package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type UserDAO interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type userDAOImpl struct {
	conn *pgx.Conn
}

func (dao *userDAOImpl) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := dao.conn.QueryRow(ctx, "SELECT id, email, password FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email, &user.PasswordHash)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, ErrUserNotFound
	case err != nil:
		return nil, err
	default:
		return &user, nil
	}
}
