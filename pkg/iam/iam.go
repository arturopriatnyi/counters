package iam

import (
	"errors"
	"net/mail"

	"github.com/google/uuid"
)

var ErrInvalidEmail = errors.New("email is not valid")

type User struct {
	ID    string
	Email string

	tokens []Token
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

func (u *User) SetToken(token Token) {
	for i := range u.tokens {
		if u.tokens[i].Issuer == token.Issuer {
			u.tokens[i] = token
			return
		}
	}

	u.tokens = append(u.tokens, token)
}

type Token struct {
	Access string
	_      string // refresh token is not used for now
	Issuer TokenIssuer
}

type TokenIssuer uint8
