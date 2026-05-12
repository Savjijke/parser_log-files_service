package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DataBaseUrl string `yaml:"data_base_url" env:"DATABASE_URL" env-default:""`
	Port        int    `yaml:"port" env:"PORT" env-default:"8080"`
	LogLevel    string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
}

func MustLoad(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}
	return cfg
}
