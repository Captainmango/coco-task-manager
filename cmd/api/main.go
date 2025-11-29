package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type (
	app struct {
		logger *slog.Logger
	}
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	a := &app{
		logger: logger,
	}

	a.logger.Info("booting application")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(customLoggingMiddleware(logger))
    
	a.logger.Info("initialising routes")
	apiV1Router := a.apiV1Router()
    r.Mount("/api/v1", apiV1Router)

	a.logger.Info("app starting")
	http.ListenAndServe(":3000", r)
}
