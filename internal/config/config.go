package config

type Mode string

const (
	ProdMode Mode = "prod"
	DevMode  Mode = "dev"
)

type Config struct {
	Mode         `env:"MODE,default=prod"`
	HTTPServer   `env:",prefix=HTTP_SERVER_"`
	GoogleOAuth2 OAuth2 `env:",prefix=GOOGLE_OAUTH2_"`
	GitHubOAuth2 OAuth2 `env:",prefix=GITHUB_OAUTH2_"`
}

type HTTPServer struct {
	Addr string `env:"ADDR,default=0.0.0.0:10000"`
}

type OAuth2 struct {
	ClientID     string   `env:"CLIENT_ID"`
	ClientSecret string   `env:"CLIENT_SECRET"`
	RedirectURL  string   `env:"REDIRECT_URL"`
	Scopes       []string `env:"SCOPES"`
}
