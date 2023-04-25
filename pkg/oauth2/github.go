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
	"golang.org/x/oauth2/github"
)

type GitHubClient struct {
	cfg    oauth2.Config
	state  string
	client http.Client
}

func NewGitHubClient(cfg Config, client http.Client) *GitHubClient {
	return &GitHubClient{
		cfg: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint:     github.Endpoint,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
		},
		state:  uuid.NewString(),
		client: client,
	}
}

func (c *GitHubClient) AuthURL() string {
	return c.cfg.AuthCodeURL(c.state)
}

func (c *GitHubClient) Exchange(ctx context.Context, state, code string) (Token, error) {
	if state != c.state {
		return Token{}, ErrInvalidState
	}

	token, err := c.cfg.Exchange(ctx, code)
	if err != nil {
		return Token{}, ErrInvalidCode
	}

	return Token{
		Access:   token.AccessToken,
		Provider: GitHub,
	}, nil
}

const gitHubUserInfoEndpoint = "https://api.github.com/user"

func (c *GitHubClient) UserInfo(ctx context.Context, token Token) (UserInfo, error) {
	req, err := nethttp.NewRequestWithContext(ctx, nethttp.MethodGet, gitHubUserInfoEndpoint, nil)
	if err != nil {
		return UserInfo{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.Access))

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
