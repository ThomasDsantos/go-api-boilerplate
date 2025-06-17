package database

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // used for migration
	_ "github.com/golang-migrate/migrate/v4/source/file"       // used for migration
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"backend/internal/config"
)

func Connect(ctx context.Context, migrations fs.FS) (*pgxpool.Pool, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("could not extract database config: %w", err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, config.PgxConfig)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	log.Debug().Msg("Connected to database, Running migrations")

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("create source: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, config.DBUrl)
	if err != nil {
		return nil, fmt.Errorf("migrate new: %w", err)
	}

	if err2 := migrator.Up(); err2 != nil && !errors.Is(err2, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to migrate database: %w", err2)
	}

	return conn, nil
}

func Health(db *pgxpool.Pool) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Error().Msgf("db down: %v", err)
		return stats, err
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := db.Stat()

	stats["open_connections"] = strconv.FormatInt(int64(dbStats.TotalConns()), 10)
	stats["idle"] = strconv.FormatInt(int64(dbStats.IdleConns()), 10)

	if dbStats.TotalConns() > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	return stats, nil
}
