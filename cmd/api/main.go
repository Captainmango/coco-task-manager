package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	coco_http "github.com/captainmango/coco-cron-parser/internal/http"
)

func main() {
	srv := coco_http.CreateApp()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Printf("application running at %s \n", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
            os.Exit(1)
            return
		}
	}()

	<-ctx.Done()
	slog.Info("application shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error(err.Error())
        os.Exit(0)
        return
	}
}
