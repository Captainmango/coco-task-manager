package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *app) attachRoutes(r *chi.Mux) *chi.Mux {
	r.Get("/api/v1/livez", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	return r
}
