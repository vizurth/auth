package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/vizurth/auth/internal/postgres"
)

type Config struct {
	Postgres postgres.Config `yaml:"postgres"`
	Port     string          `yaml:"port" env:"PORT" env-default:"50051"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("./config/config.yaml", &cfg); err != nil {
		fmt.Println(err)
		if err = cleanenv.ReadEnv(&cfg); err != nil {
			return Config{}, err
		}
	}
	return cfg, nil
}
