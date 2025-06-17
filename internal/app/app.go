package app

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"backend/internal/config"
	"backend/internal/database"
	store_pkg "backend/internal/store"
)

type App struct {
	Config config.ServerConfig
	files  fs.FS
	API    *http.Server
	DB     *pgxpool.Pool
	Store  *store_pkg.Queries
}

func New(files fs.FS) *App {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	app := App{
		Config: *cfg,
		files:  files,
	}
	app.CreateLogger()

	return &app
}

func (a *App) CreateLogger() {
	zerolog.SetGlobalLevel(a.Config.LogLevel)
	if slices.Contains([]string{"development", "local"}, a.Config.Environment) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
	log.Info().Str("environment", a.Config.Environment).Msg("Logger is setup")
}

func (a *App) Start(ctx context.Context) error {
	db, err := database.Connect(ctx, a.files)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	a.DB = db
	a.Store = store_pkg.New(db)
	a.createServer()

	errCh := make(chan error, 1)
	go func() {
		err2 := a.API.ListenAndServe()
		if err2 != nil && !errors.Is(err2, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to listen and serve: %w", err2)
		}
		close(errCh)
	}()
	log.Info().Msgf("Server running on port %v", a.Config.Port)

	select {
	// Wait until we receive SIGINT (ctrl+c on cli)
	case <-ctx.Done():
		break
	case err3 := <-errCh:
		return err3
	}

	sCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	return a.API.Shutdown(sCtx)
}
