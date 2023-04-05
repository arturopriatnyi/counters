package iam

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewManager(t *testing.T) {
	wantManager := &Manager{users: NewMockUserStorage(gomock.NewController(t))}

	m := NewManager(wantManager.users)

	if !reflect.DeepEqual(m, wantManager) {
		t.Errorf("want: %+v, got: %+v", wantManager, m)
	}
}

func TestManager_SignInWithOAuth2(t *testing.T) {
	errUnexpected := errors.New("unexpected error")

	for name, tt := range map[string]struct {
		users   func(*gomock.Controller) UserStorage
		email   string
		token   Token
		wantErr error
	}{
		"NewUserIsSignedIn": {
			users: func(c *gomock.Controller) UserStorage {
				s := NewMockUserStorage(c)

				s.
					EXPECT().
					Get("x@x.x").
					Return(
						nil,
						ErrUserNotFound,
					)
				s.
					EXPECT().
					Set(gomock.AssignableToTypeOf(&User{})).
					Return(nil)

				return s
			},
			email:   "x@x.x",
			token:   Token{Access: "accessToken", Issuer: 1},
			wantErr: nil,
		},
		"ExistingUserIsSignedIn": {
			users: func(c *gomock.Controller) UserStorage {
				s := NewMockUserStorage(c)

				s.
					EXPECT().
					Get("x@x.x").
					Return(
						&User{
							ID:    "x-x-x-x-x",
							Email: "x@x.x",
							tokens: []Token{
								{Access: "accessToken1", Issuer: 1},
							},
						},
						nil,
					)
				s.
					EXPECT().
					Set(
						&User{
							ID:    "x-x-x-x-x",
							Email: "x@x.x",
							tokens: []Token{
								{Access: "accessToken1", Issuer: 1},
								{Access: "accessToken2", Issuer: 2},
							},
						},
					).
					Return(nil)

				return s
			},
			email:   "x@x.x",
			token:   Token{Access: "accessToken2", Issuer: 2},
			wantErr: nil,
		},
		"GetUserUnexpectedError": {
			users: func(c *gomock.Controller) UserStorage {
				s := NewMockUserStorage(c)

				s.
					EXPECT().
					Get("x@x.x").
					Return(
						nil,
						errUnexpected,
					)

				return s
			},
			email:   "x@x.x",
			token:   Token{Access: "accessToken1", Issuer: 1},
			wantErr: errUnexpected,
		},
		"SetUserUnexpectedError": {
			users: func(c *gomock.Controller) UserStorage {
				s := NewMockUserStorage(c)

				s.
					EXPECT().
					Get("x@x.x").
					Return(
						nil,
						ErrUserNotFound,
					)
				s.
					EXPECT().
					Set(gomock.AssignableToTypeOf(&User{})).
					Return(errUnexpected)

				return s
			},
			email:   "x@x.x",
			token:   Token{Access: "accessToken", Issuer: 1},
			wantErr: errUnexpected,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := &Manager{users: tt.users(gomock.NewController(t))}

			err := m.SignInWithOAuth2(context.TODO(), tt.email, tt.token)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
