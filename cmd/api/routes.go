package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *app) apiV1Router(r *chi.Mux) *chi.Mux {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/livez", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		})
	})


	return r
}
