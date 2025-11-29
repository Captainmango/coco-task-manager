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
    a.logger.Info("initialising routes")
	r.Use(customLoggingMiddleware(logger))

	router := a.attachRoutes(r)

    a.logger.Info("app starting")
	http.ListenAndServe(":3000", router)
}
