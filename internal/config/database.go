package config

import (
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrMissingDatabaseURL = errors.New("DATABASE_URL env missing")
)

type DatabaseConfig struct {
	PgxConfig *pgxpool.Config
	DBUrl     string
}

func LoadConfig() (*DatabaseConfig, error) {
	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, ErrMissingDatabaseURL
	}

	pgxPoolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	return &DatabaseConfig{
		PgxConfig: pgxPoolConfig,
		DBUrl:     dbURL,
	}, nil
}
