package config

import "github.com/caarlos0/env"

type Config struct {
	HttpPort  string `env:"HTTP_PORT" envDefault:"3000"`
	DbConnUrl string `env:"DB_CONN_URL" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
}

func MustLoad() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
