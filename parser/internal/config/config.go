package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DataBaseUrl string `env:"DATABASE_URL" env-default:"postgres://postgres:postgres@db:5432/parser?sslmode=disable"`
	Port        int    `env:"PORT" env-default:"8080"`
	LogLevel    string `env:"LOG_LEVEL" env-default:"DEBUG"`
}

func MustLoad() Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config %q:", err)
	}
	return cfg
}
