package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *app) apiV1Router() *chi.Mux {
	r := chi.NewRouter()
	
	r.Get("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	return r
}
