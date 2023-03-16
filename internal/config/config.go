package config

type Mode string

const (
	ProdMode Mode = "prod"
	DevMode  Mode = "dev"
)

type Config struct {
	Mode       `env:"MODE,default=prod"`
	HTTPServer `env:",prefix=HTTP_SERVER_"`
}

type HTTPServer struct {
	Addr string `env:"ADDR,default=0.0.0.0:10000"`
}
