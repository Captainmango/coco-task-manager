package coco_http

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	coco_cli "github.com/captainmango/coco-cron-parser/internal/cli"
	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/resources"
)

type (
	app struct {
		logger           *slog.Logger
		resources        resources.Resources
		commandsRegistry coco_cli.CommandFinder
	}
)

func CreateApp() *http.Server {
	config.BootstrapConfig(
		config.WithDotEnv(),
	)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	a := &app{
		logger:           logger,
		resources:        resources.CreateResources(),
		commandsRegistry: coco_cli.CommandRegistry,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(customLoggingMiddleware(logger))

	return &http.Server{
		Addr:        ":3000",
		Handler:     a.apiV1Router(r),
		ReadTimeout: 5 * time.Second,
	}
}
