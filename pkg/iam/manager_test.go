package iam

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"reflect"
	"testing"

	"counters/pkg/oauth2"

	"github.com/golang/mock/gomock"
)

func TestNewManager(t *testing.T) {
	wantManager := &Manager{
		users:  NewMockUserStorage(gomock.NewController(t)),
		oauth2: map[oauth2.Provider]oauth2.Client{},
	}

	m := NewManager(wantManager.users, wantManager.oauth2)

	if !reflect.DeepEqual(m, wantManager) {
		t.Errorf("want: %+v, got: %+v", wantManager, m)
	}
}

func TestManager_OAuth2URL(t *testing.T) {
	for name, tt := range map[string]struct {
		provider      oauth2.Provider
		oauth2        func(*gomock.Controller) map[oauth2.Provider]oauth2.Client
		wantOAuth2URL string
		wantErr       error
	}{
		"OK": {
			provider: oauth2.Google,
			oauth2: func(c *gomock.Controller) map[oauth2.Provider]oauth2.Client {
				client := oauth2.NewMockClient(c)

				client.
					EXPECT().
					AuthURL().
					Return("oauth2.url")

				return map[oauth2.Provider]oauth2.Client{
					oauth2.Google: client,
				}
			},
			wantOAuth2URL: "oauth2.url",
			wantErr:       nil,
		},
		"ErrInvalidOAuth2Provider": {
			provider: oauth2.GitHub,
			oauth2: func(c *gomock.Controller) map[oauth2.Provider]oauth2.Client {
				return map[oauth2.Provider]oauth2.Client{}
			},
			wantOAuth2URL: "",
			wantErr:       ErrInvalidOAuth2Provider,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := &Manager{oauth2: tt.oauth2(gomock.NewController(t))}

			oauth2URL, err := m.OAuth2URL(tt.provider)

			if oauth2URL != tt.wantOAuth2URL {
				t.Errorf("want: %s, got: %s", tt.wantOAuth2URL, oauth2URL)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestManager_SignInWithOAuth2(t *testing.T) {
	errUnexpected := errors.New("unexpected error")

	for name, tt := range map[string]struct {
		provider    oauth2.Provider
		state, code string
		users       func(*gomock.Controller) UserStorage
		oauth2      func(controller *gomock.Controller) map[oauth2.Provider]oauth2.Client
		wantToken   oauth2.Token
		wantErr     error
	}{
		"OK_NewUserIsSignedIn": {
			provider: oauth2.Google,
			state:    "state",
			code:     "code",
			users: func(c *gomock.Controller) UserStorage {
				s := NewMockUserStorage(c)

				s.
					EXPECT().
					Get("x@x.x").
					Return(nil, ErrUserNotFound)
				s.
					EXPECT().
					Set(gomock.AssignableToTypeOf(&User{})).
					Return(nil)

				return s
			},
			oauth2: func(c *gomock.Controller) map[oauth2.Provider]oauth2.Client {
				client := oauth2.NewMockClient(c)

				client.
					EXPECT().
					Exchange(context.TODO(), "state", "code").
					Return(
						oauth2.Token{
							Access:   "accessToken",
							Provider: oauth2.Google,
						},
						nil,
					)
				client.
					EXPECT().
					UserInfo(
						context.TODO(),
						oauth2.Token{
							Access:   "accessToken",
							Provider: oauth2.Google,
						},
					).
					Return(oauth2.UserInfo{Email: "x@x.x"}, nil)

				return map[oauth2.Provider]oauth2.Client{
					oauth2.Google: client,
				}
			},
			wantToken: oauth2.Token{
				Access:   "accessToken",
				Provider: oauth2.Google,
			},
			wantErr: nil,
		},
		"OK_ExistingUserIsSignedIn": {
			provider: oauth2.Google,
			state:    "state",
			code:     "code",
			users: func(c *gomock.Controller) UserStorage {
				s := NewMockUserStorage(c)

				s.
					EXPECT().
					Get("x@x.x").
					Return(&User{
						ID:    uuid.NewString(),
						Email: "x@x.x",
						tokens: []oauth2.Token{
							{Access: "accessToken", Provider: oauth2.GitHub},
						},
					}, nil)
				s.
					EXPECT().
					Set(gomock.AssignableToTypeOf(&User{})).
					Return(nil)

				return s
			},
			oauth2: func(c *gomock.Controller) map[oauth2.Provider]oauth2.Client {
				client := oauth2.NewMockClient(c)

				client.
					EXPECT().
					Exchange(context.TODO(), "state", "code").
					Return(
						oauth2.Token{
							Access:   "newAccessToken",
							Provider: oauth2.Google,
						},
						nil,
					)
				client.
					EXPECT().
					UserInfo(
						context.TODO(),
						oauth2.Token{
							Access:   "newAccessToken",
							Provider: oauth2.Google,
						},
					).
					Return(oauth2.UserInfo{Email: "x@x.x"}, nil)

				return map[oauth2.Provider]oauth2.Client{
					oauth2.Google: client,
				}
			},
			wantToken: oauth2.Token{
				Access:   "newAccessToken",
				Provider: oauth2.Google,
			},
			wantErr: nil,
		},
		"ErrInvalidOAuth2Provider": {
			provider: oauth2.GitHub,
			state:    "state",
			code:     "code",
			users: func(c *gomock.Controller) UserStorage {
				return NewMockUserStorage(c)
			},
			oauth2: func(c *gomock.Controller) map[oauth2.Provider]oauth2.Client {
				return map[oauth2.Provider]oauth2.Client{
					oauth2.Google: oauth2.NewMockClient(c),
				}
			},
			wantToken: oauth2.Token{},
			wantErr:   ErrInvalidOAuth2Provider,
		},
		"ExchangeUnexpectedError": {
			provider: oauth2.Google,
			state:    "state",
			code:     "code",
			users: func(c *gomock.Controller) UserStorage {
				return NewMockUserStorage(c)
			},
			oauth2: func(c *gomock.Controller) map[oauth2.Provider]oauth2.Client {
				client := oauth2.NewMockClient(c)

				client.
					EXPECT().
					Exchange(context.TODO(), "state", "code").
					Return(oauth2.Token{}, errUnexpected)

				return map[oauth2.Provider]oauth2.Client{
					oauth2.Google: client,
				}
			},
			wantToken: oauth2.Token{},
			wantErr:   errUnexpected,
		},
		"UserInfoUnexpectedError": {
			provider: oauth2.Google,
			state:    "state",
			code:     "code",
			users: func(c *gomock.Controller) UserStorage {
				return NewMockUserStorage(c)
			},
			oauth2: func(c *gomock.Controller) map[oauth2.Provider]oauth2.Client {
				client := oauth2.NewMockClient(c)

				client.
					EXPECT().
					Exchange(context.TODO(), "state", "code").
					Return(
						oauth2.Token{
							Access:   "accessToken",
							Provider: oauth2.Google,
						}, nil,
					)
				client.
					EXPECT().
					UserInfo(
						context.TODO(),
						oauth2.Token{
							Access:   "accessToken",
							Provider: oauth2.Google,
						}).
					Return(oauth2.UserInfo{}, errUnexpected)

				return map[oauth2.Provider]oauth2.Client{
					oauth2.Google: client,
				}
			},
			wantToken: oauth2.Token{},
			wantErr:   errUnexpected,
		},
		"GetUserUnexptedError": {
			provider: oauth2.Google,
			state:    "state",
			code:     "code",
			users: func(c *gomock.Controller) UserStorage {
				s := NewMockUserStorage(c)

				s.
					EXPECT().
					Get("x@x.x").
					Return(nil, errUnexpected)

				return s
			},
			oauth2: func(c *gomock.Controller) map[oauth2.Provider]oauth2.Client {
				client := oauth2.NewMockClient(c)

				client.
					EXPECT().
					Exchange(context.TODO(), "state", "code").
					Return(
						oauth2.Token{
							Access:   "accessToken",
							Provider: oauth2.Google,
						},
						nil,
					)
				client.
					EXPECT().
					UserInfo(
						context.TODO(),
						oauth2.Token{
							Access:   "accessToken",
							Provider: oauth2.Google,
						},
					).
					Return(oauth2.UserInfo{Email: "x@x.x"}, nil)

				return map[oauth2.Provider]oauth2.Client{
					oauth2.Google: client,
				}
			},
			wantToken: oauth2.Token{},
			wantErr:   errUnexpected,
		},
		"NewUserInvalidEmailError": {
			provider: oauth2.Google,
			state:    "state",
			code:     "code",
			users: func(c *gomock.Controller) UserStorage {
				s := NewMockUserStorage(c)

				s.
					EXPECT().
					Get("x").
					Return(nil, ErrUserNotFound)

				return s
			},
			oauth2: func(c *gomock.Controller) map[oauth2.Provider]oauth2.Client {
				client := oauth2.NewMockClient(c)

				client.
					EXPECT().
					Exchange(context.TODO(), "state", "code").
					Return(
						oauth2.Token{
							Access:   "accessToken",
							Provider: oauth2.Google,
						},
						nil,
					)
				client.
					EXPECT().
					UserInfo(
						context.TODO(),
						oauth2.Token{
							Access:   "accessToken",
							Provider: oauth2.Google,
						},
					).
					Return(oauth2.UserInfo{Email: "x"}, nil)

				return map[oauth2.Provider]oauth2.Client{
					oauth2.Google: client,
				}
			},
			wantToken: oauth2.Token{},
			wantErr:   ErrInvalidEmail,
		},
		"SetUserUnexpectedError": {
			provider: oauth2.Google,
			state:    "state",
			code:     "code",
			users: func(c *gomock.Controller) UserStorage {
				s := NewMockUserStorage(c)

				s.
					EXPECT().
					Get("x@x.x").
					Return(nil, ErrUserNotFound)
				s.
					EXPECT().
					Set(gomock.AssignableToTypeOf(&User{})).
					Return(errUnexpected)

				return s
			},
			oauth2: func(c *gomock.Controller) map[oauth2.Provider]oauth2.Client {
				client := oauth2.NewMockClient(c)

				client.
					EXPECT().
					Exchange(context.TODO(), "state", "code").
					Return(
						oauth2.Token{
							Access:   "accessToken",
							Provider: oauth2.Google,
						},
						nil,
					)
				client.
					EXPECT().
					UserInfo(
						context.TODO(),
						oauth2.Token{
							Access:   "accessToken",
							Provider: oauth2.Google,
						},
					).
					Return(oauth2.UserInfo{Email: "x@x.x"}, nil)

				return map[oauth2.Provider]oauth2.Client{
					oauth2.Google: client,
				}
			},
			wantToken: oauth2.Token{},
			wantErr:   errUnexpected,
		},
	} {
		t.Run(name, func(t *testing.T) {
			c := gomock.NewController(t)
			m := &Manager{
				users:  tt.users(c),
				oauth2: tt.oauth2(c),
			}

			token, err := m.SignInWithOAuth2(context.TODO(), tt.provider, tt.state, tt.code)

			if token != tt.wantToken {
				t.Errorf("want: %+v, got: %+v", tt.wantToken, token)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
