package iam

import (
	"errors"
	"reflect"
	"testing"

	"counters/pkg/oauth2"
)

func TestNewUser(t *testing.T) {
	for name, tt := range map[string]struct {
		email   string
		wantErr error
	}{
		"OK": {
			email:   "x@x.x",
			wantErr: nil,
		},
		"ErrInvalidEmail": {
			email:   "",
			wantErr: ErrInvalidEmail,
		},
	} {
		t.Run(name, func(t *testing.T) {
			u, err := NewUser(tt.email)

			if err == nil {
				if u.ID == "" {
					t.Error("want: <any UUID>, got: ")
				}
				if u.Email != tt.email {
					t.Errorf("want: %s, got: %s", tt.email, u.Email)
				}
				if u.tokens != nil {
					t.Errorf("want: [], got: %v", u.tokens)
				}
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestUser_SetToken(t *testing.T) {
	for name, tt := range map[string]struct {
		user       User
		token      oauth2.Token
		wantTokens []oauth2.Token
	}{
		"FirstTokenIsSet": {
			user:  User{tokens: nil},
			token: oauth2.Token{Access: "accessToken", Provider: 1},
			wantTokens: []oauth2.Token{
				{Access: "accessToken", Provider: 1},
			},
		},
		"SecondTokenIsSet": {
			user: User{
				tokens: []oauth2.Token{
					{Access: "accessToken1", Provider: 1},
				},
			},
			token: oauth2.Token{Access: "accessToken2", Provider: 2},
			wantTokens: []oauth2.Token{
				{Access: "accessToken1", Provider: 1},
				{Access: "accessToken2", Provider: 2},
			},
		},
		"TokenIsOverridden": {
			user: User{
				tokens: []oauth2.Token{
					{Access: "accessToken", Provider: 1},
				},
			},
			token: oauth2.Token{Access: "newAccessToken", Provider: 1},
			wantTokens: []oauth2.Token{
				{Access: "newAccessToken", Provider: 1},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			tt.user.SetToken(tt.token)

			if !reflect.DeepEqual(tt.user.tokens, tt.wantTokens) {
				t.Errorf("want: %+v, got: %+v", tt.wantTokens, tt.user.tokens)
			}
		})
	}
}
