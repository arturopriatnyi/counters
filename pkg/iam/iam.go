package iam

import (
	"errors"
	"net/mail"

	"counters/pkg/oauth2"

	"github.com/google/uuid"
)

var ErrInvalidEmail = errors.New("email is not valid")

type User struct {
	ID    string
	Email string

	tokens []oauth2.Token
}

func NewUser(email string) (*User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrInvalidEmail
	}

	return &User{
		ID:    uuid.NewString(),
		Email: email,
	}, nil
}

func (u *User) SetToken(token oauth2.Token) {
	for i := range u.tokens {
		if u.tokens[i].Provider == token.Provider {
			u.tokens[i] = token
			return
		}
	}

	u.tokens = append(u.tokens, token)
}
