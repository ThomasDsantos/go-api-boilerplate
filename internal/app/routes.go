package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"backend/internal/app/handlers"
	"backend/internal/app/middleware"
)

func (a *App) loadRoutes(api huma.API) {
	greetingHandler := handlers.NewGreetingHandler(a.Store)

	huma.Register(api, huma.Operation{
		OperationID: "get-greeting",
		Method:      http.MethodGet,
		Path:        "/greeting/{name}",
		Summary:     "Get a greeting",
	}, greetingHandler.GetGreeting)
}

func (a *App) getApiRouter() *chi.Mux {
	apiSubRouter := chi.NewRouter()
	humaCfg := huma.DefaultConfig(a.Config.ServiceName, "1.0.0")
	humaCfg.Info.Contact = &huma.Contact{Email: "dev@example.com"}
	humaCfg.Info.Description = "This is a sample API built with Huma and Chi."
	humaCfg.DocsPath = "/docs"
	humaCfg.Servers = []*huma.Server{
		{
			URL:         a.Config.APIBasePath,
			Description: "local API Server (" + a.Config.ServiceName + ")",
		},
	}
	humaAPI := humachi.New(apiSubRouter, humaCfg)
	a.loadRoutes(humaAPI)
	return apiSubRouter
}

func (a *App) createServer() (*http.Server, error) {
	router := chi.NewRouter()
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Recoverer)

	router.Use(hlog.NewHandler(log.Logger))
	router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status_code", status).
			Dur("duration_ms", duration).
			Msg("Request")
	}))
	router.Use(hlog.RemoteAddrHandler("ip"))
	router.Use(middleware.HumaMiddleware)

	apiRouter := a.getApiRouter()
	router.Mount(a.Config.APIBasePath, apiRouter)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%v", a.Config.Port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return httpServer, nil
}
