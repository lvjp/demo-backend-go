package auth

import (
	"errors"
)

var ErrAccessDenied = errors.New("access denied")

type SessionService interface {
	Create(input SessionCreateInput) (err error)
}

func NewSessionService() SessionService {
	return &impl{}
}

type impl struct{}

func (*impl) Create(input SessionCreateInput) error {
	if input.ID != "me@example.com" || input.Password != "pa$$w0rd" {
		return ErrAccessDenied
	}

	return nil
}
