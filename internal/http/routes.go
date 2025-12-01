package coco_http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *app) apiV1Router(r *chi.Mux) *chi.Mux {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/livez", a.handleLivez)

		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", a.handleGetTasks)
			r.Get("/scheduled", a.handleGetScheduledTasks)

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {

			})
		})
	})

	return r
}
