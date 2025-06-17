package main

import (
	"context"
	"embed"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"

	"backend/internal/app"
)

//go:embed migrations/*.sql
var files embed.FS

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := app.New(files)

	if err := app.Start(ctx); err != nil {
		log.Err(err).Msgf("failed to start app")
	}
}
