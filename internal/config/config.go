package config

type Config struct {
	HTTPServer `env:",prefix=HTTP_SERVER_"`
}

type HTTPServer struct {
	Addr string `env:"ADDR,default=0.0.0.0:10000"`
}
