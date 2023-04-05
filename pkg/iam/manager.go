package iam

import "context"

type Manager struct {
	users UserStorage
}

func NewManager(users UserStorage) *Manager {
	return &Manager{users: users}
}

func (m *Manager) SignInWithOAuth2(_ context.Context, email string, token Token) error {
	u, err := m.users.Get(email)
	if err != nil && err != ErrUserNotFound {
		return err
	}

	if err == ErrUserNotFound {
		if u, err = NewUser(email); err != nil {
			return err
		}
	}

	u.SetToken(token)

	return m.users.Set(u)
}
