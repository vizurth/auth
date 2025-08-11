package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"host" envconfig:"PG_HOST"`
	Port     string `yaml:"port" envconfig:"PG_PORT"`
	Username string `yaml:"username" envconfig:"PG_USERNAME"`
	Password string `yaml:"password" envconfig:"PG_PASSWORD"`
	Database string `yaml:"database" envconfig:"PG_DB"`
}

func New(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connString := cfg.GetConnString()

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Config) GetConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}
