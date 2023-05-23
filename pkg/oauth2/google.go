package oauth2

import (
	"context"
	"encoding/json"
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

func (c *GoogleClient) AuthURL() string {
	return c.cfg.AuthCodeURL(c.state)
}

func (c *GoogleClient) Exchange(ctx context.Context, state, code string) (Token, error) {
	if state != c.state {
		return Token{}, ErrInvalidState
	}

	token, err := c.cfg.Exchange(ctx, code)
	if err != nil {
		return Token{}, ErrInvalidCode
	}

	return Token{
		Access:   token.AccessToken,
		Provider: Google,
	}, nil
}

const googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v2/userinfo"

func (c *GoogleClient) UserInfo(ctx context.Context, token Token) (UserInfo, error) {
	url := fmt.Sprintf("%s?access_token=%s", googleUserInfoEndpoint, token.Access)

	req, err := nethttp.NewRequestWithContext(ctx, nethttp.MethodGet, url, nil)
	if err != nil {
		return UserInfo{}, err
	}

	res, err := c.client.Do(req)
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
