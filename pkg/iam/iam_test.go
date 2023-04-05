package iam

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	for name, tt := range map[string]struct {
		email   string
		wantErr error
	}{
		"UserIsCreated": {
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
		tokens     []Token
		token      Token
		wantTokens []Token
	}{
		"FirstTokenIsSet": {
			tokens: nil,
			token:  Token{Access: "accessToken", Issuer: 1},
			wantTokens: []Token{
				{Access: "accessToken", Issuer: 1},
			},
		},
		"SecondTokenIsSet": {
			tokens: []Token{
				{Access: "accessToken1", Issuer: 1},
			},
			token: Token{Access: "accessToken2", Issuer: 2},
			wantTokens: []Token{
				{Access: "accessToken1", Issuer: 1},
				{Access: "accessToken2", Issuer: 2},
			},
		},
		"TokenIsOverridden": {
			tokens: []Token{
				{Access: "accessToken", Issuer: 1},
			},
			token: Token{Access: "newAccessToken", Issuer: 1},
			wantTokens: []Token{
				{Access: "newAccessToken", Issuer: 1},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			u := &User{tokens: tt.tokens}

			u.SetToken(tt.token)

			if !reflect.DeepEqual(u.tokens, tt.wantTokens) {
				t.Errorf("want: %+v, got: %+v", tt.wantTokens, u.tokens)
			}
		})
	}
}
