package config

import (
	"os"
)

type (
	Config struct {
		Server Server
	}

	Server struct {
		Port string `env:"PORT,required,notEmpty"`
	}
)

const (
	DefaultPort = "8080"
)

func New() *Config {
	cfg := &Config{}
	cfg.Server.Port = os.Getenv("PORT")
	if len(cfg.Server.Port) == 0 {
		cfg.Server.Port = DefaultPort
	}

	return cfg
}
