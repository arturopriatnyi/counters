//go:generate mockgen -source=storage.go -destination=mock.go -package=iam
package iam

import (
	"errors"
	"sync"
)

var ErrUserNotFound = errors.New("user not found")

type UserStorage interface {
	Set(user *User) error
	Get(email string) (*User, error)
}

type UserMemoryStorage struct {
	mu    sync.RWMutex
	users map[string]*User
}

func NewUserMemoryStorage() *UserMemoryStorage {
	return &UserMemoryStorage{users: make(map[string]*User)}
}

func (s *UserMemoryStorage) Set(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.Email] = user

	return nil
}

func (s *UserMemoryStorage) Get(email string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[email]
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}
