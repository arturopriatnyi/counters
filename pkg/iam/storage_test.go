package iam

import (
	"errors"
	"reflect"
	"testing"

	"counters/pkg/oauth2"
)

func TestNewUserMemoryStorage(t *testing.T) {
	wantStorage := &UserMemoryStorage{users: make(map[string]*User)}

	storage := NewUserMemoryStorage()

	if !reflect.DeepEqual(storage, wantStorage) {
		t.Errorf("want: %+v, got: %+v", wantStorage, storage)
	}
}

func TestUserMemoryStorage_Set(t *testing.T) {
	storage := &UserMemoryStorage{users: make(map[string]*User)}
	user := &User{ID: "x-x-x-x-x", Email: "x@x.x", tokens: []oauth2.Token{}}

	err := storage.Set(user)

	if err != nil {
		t.Errorf("want: <nil>, got: %v", err)
	}
	if u, ok := storage.users[user.Email]; !ok || !reflect.DeepEqual(u, user) {
		t.Errorf("want: %+v, got: %+v", user, u)
	}
}

func TestUserMemoryStorage_Get(t *testing.T) {
	for name, tt := range map[string]struct {
		storage  *UserMemoryStorage
		email    string
		wantUser *User
		wantErr  error
	}{
		"OK": {
			storage: &UserMemoryStorage{
				users: map[string]*User{
					"x@x.x": {Email: "x@x.x"},
				},
			},
			email:    "x@x.x",
			wantUser: &User{Email: "x@x.x"},
			wantErr:  nil,
		},
		"ErrUserNotFound": {
			storage:  &UserMemoryStorage{users: map[string]*User{}},
			email:    "x@x.x",
			wantUser: nil,
			wantErr:  ErrUserNotFound,
		},
	} {
		t.Run(name, func(t *testing.T) {
			u, err := tt.storage.Get(tt.email)

			if !reflect.DeepEqual(u, tt.wantUser) {
				t.Errorf("want: %+v, got: %+v", tt.wantUser, u)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
