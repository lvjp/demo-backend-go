package auth

import (
	"context"
	"errors"

	"go.lvjp.me/demo-backend-go/internal/app/db"
	"go.lvjp.me/demo-backend-go/pkg/hashutils/password"
)

var ErrAccessDenied = errors.New("access denied")

type SessionService interface {
	Create(ctx context.Context, input SessionCreateInput) (err error)
}

func NewSessionService(conn db.Connector) SessionService {
	return &impl{
		userDAO: conn.UserDAO(),
	}
}

type impl struct {
	userDAO db.UserDAO
}

func (impl *impl) Create(ctx context.Context, input SessionCreateInput) error {
	user, err := impl.userDAO.GetByEmail(ctx, input.Email)
	if errors.Is(err, db.ErrUserNotFound) {
		return ErrAccessDenied
	}
	if err != nil {
		return err
	}

	isSame, err := password.IsSame([]byte(input.Password), user.PasswordHash)
	if err != nil {
		return err
	}
	if !isSame {
		return ErrAccessDenied
	}

	return nil
}
