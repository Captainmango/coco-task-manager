package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func customLoggingMiddleware(l *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqStart := time.Now()

			rw := &responseWriter{w, 200}
			next.ServeHTTP(rw, r)

			duration := time.Since(reqStart)
			requestId := middleware.GetReqID(r.Context())

			l.InfoContext(r.Context(), "request completed", 
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
                slog.Int("status", rw.statusCode),
                slog.Duration("duration", duration),
                slog.String("remote_addr", r.RemoteAddr),
				slog.String("request_id", requestId),
			)
		})
	}
}