package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	nethttp "net/http"

	"counters/pkg/http"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleClient struct {
	cfg    oauth2.Config
	state  string
	client http.Client
}

func NewGoogleClient(cfg Config, client http.Client) *GoogleClient {
	return &GoogleClient{
		cfg: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
		},
		state:  uuid.NewString(),
		client: client,
	}
}

func (g *GoogleClient) AuthURL() string {
	return g.cfg.AuthCodeURL(g.state)
}

var (
	ErrInvalidState = errors.New("invalid state")
	ErrInvalidCode  = errors.New("invalid code")
)

func (g *GoogleClient) Exchange(ctx context.Context, state, code string) (Token, error) {
	if state != g.state {
		return Token{}, ErrInvalidState
	}

	token, err := g.cfg.Exchange(ctx, code)
	if err != nil {
		return Token{}, ErrInvalidCode
	}

	return Token{
		Access:   token.AccessToken,
		Provider: Google,
	}, nil
}

const googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v2/userinfo"

func (g *GoogleClient) UserInfo(ctx context.Context, t Token) (UserInfo, error) {
	url := fmt.Sprintf("%s?access_token=%s", googleUserInfoEndpoint, t.Access)

	req, err := nethttp.NewRequestWithContext(ctx, nethttp.MethodGet, url, nil)
	if err != nil {
		return UserInfo{}, err
	}

	res, err := g.client.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return UserInfo{}, err
	}

	var info UserInfo
	return info, json.Unmarshal(body, &info)
}
