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
	Api    *http.Server
	Db     *pgxpool.Pool
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
	a.Db = db

	a.Store = store_pkg.New(db)

	srv, err := a.createServer()
	if err != nil {
		return fmt.Errorf("failed when creating server: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to listen and serve: %w", err)
		}
		close(errCh)
	}()
	log.Info().Msgf("Server running on port %v", a.Config.Port)

	select {
	// Wait until we receive SIGINT (ctrl+c on cli)
	case <-ctx.Done():
		break
	case err := <-errCh:
		return err
	}

	sCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	return srv.Shutdown(sCtx)
}
