//go:generate mockgen -source=client.go -destination=mock.go -package=oauth2
package oauth2

import (
	"context"
	"errors"
)

var (
	ErrInvalidState = errors.New("invalid state")
	ErrInvalidCode  = errors.New("invalid code")
)

type Client interface {
	AuthURL() string
	Exchange(ctx context.Context, state, code string) (Token, error)
	UserInfo(ctx context.Context, token Token) (UserInfo, error)
}

type Token struct {
	Access   string
	Provider Provider
}

type Provider uint8

const (
	Google Provider = iota + 1
	GitHub
)

type UserInfo struct {
	Email string `json:"email"`
}

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}
