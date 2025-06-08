package database

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rs/zerolog/log"

	"backend/internal/config"
)

var (
	ErrMissingMigrationsPath = errors.New("MIGRATIONS_PATH env missing")
	ErrMissingDatabaseURL    = errors.New("DATABASE_URL env missing")
)

func loadConfigFromURL() (*pgxpool.Config, error) {
	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, fmt.Errorf("must set DATABASE_URL env var")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, nil
}

func loadConfig() (*pgxpool.Config, error) {
	cfg, err := config.NewDatabase()
	if err != nil {
		return loadConfigFromURL()
	}

	return pgxpool.ParseConfig(fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	))
}

func dbURL() (string, error) {
	cfg, err := config.NewDatabase()
	if err != nil {
		dbURL, ok := os.LookupEnv("DATABASE_URL")
		if !ok {
			return "", fmt.Errorf("must set DATABASE_URL env var")
		}

		return dbURL, nil
	}

	return cfg.URL(), nil
}

func Connect(ctx context.Context, migrations fs.FS) (*pgxpool.Pool, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	log.Debug().Msg("Connected to database, Running migrations")

	url, err := dbURL()
	if err != nil {
		return nil, err
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("create source: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, url)
	if err != nil {
		return nil, fmt.Errorf("migrate new: %s", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return conn, nil
}

func Health(db *pgxpool.Pool) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatal().Msgf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := db.Stat()

	stats["open_connections"] = fmt.Sprint(dbStats.TotalConns())
	stats["idle"] = fmt.Sprint(dbStats.IdleConns())

	if dbStats.TotalConns() > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	return stats
}

