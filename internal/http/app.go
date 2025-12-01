package coco_http

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/resources"
	"github.com/captainmango/coco-cron-parser/internal/commands"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type (
	app struct {
		logger *slog.Logger
		resources resources.Resources
		cmdRegistry commands.RegistryContainer
	}
)

func CreateApp() *http.Server {
	config.BootstrapConfig(
		config.WithDotEnv(),
	)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	a := &app{
		logger: logger,
		resources: resources.CreateResources(),
		cmdRegistry: *commands.Registry,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(customLoggingMiddleware(logger))

    return &http.Server{
		Addr: ":3000",
		Handler: a.apiV1Router(r),
	}
}