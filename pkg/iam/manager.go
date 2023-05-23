package iam

import (
	"context"
	"errors"

	"counters/pkg/oauth2"
)

type Manager struct {
	users UserStorage

	oauth2 map[oauth2.Provider]oauth2.Client
}

func NewManager(users UserStorage, oauth2 map[oauth2.Provider]oauth2.Client) *Manager {
	return &Manager{users: users, oauth2: oauth2}
}

var ErrInvalidOAuth2Provider = errors.New("invalid OAuth2 provider")

func (m *Manager) OAuth2URL(provider oauth2.Provider) (string, error) {
	c, ok := m.oauth2[provider]
	if !ok {
		return "", ErrInvalidOAuth2Provider
	}

	return c.AuthURL(), nil
}

func (m *Manager) SignInWithOAuth2(ctx context.Context, provider oauth2.Provider, state, code string) (oauth2.Token, error) {
	c, ok := m.oauth2[provider]
	if !ok {
		return oauth2.Token{}, ErrInvalidOAuth2Provider
	}

	token, err := c.Exchange(ctx, state, code)
	if err != nil {
		return oauth2.Token{}, err
	}

	info, err := c.UserInfo(ctx, token)
	if err != nil {
		return oauth2.Token{}, err
	}

	u, err := m.users.Get(info.Email)
	if err != nil && err != ErrUserNotFound {
		return oauth2.Token{}, err
	}

	if err == ErrUserNotFound {
		if u, err = NewUser(info.Email); err != nil {
			return oauth2.Token{}, err
		}
	}

	u.SetToken(token)

	if err = m.users.Set(u); err != nil {
		return oauth2.Token{}, err
	}

	return token, nil
}
