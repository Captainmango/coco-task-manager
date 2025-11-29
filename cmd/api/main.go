package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	a.logger.Info("app starting")

    server := &http.Server{Addr: ":3000", Handler: a.apiV1Router(r)}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
            os.Exit(1)
            return
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error(err.Error())
        os.Exit(0)
        return
	}
}
