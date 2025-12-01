package coco_http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *app) apiV1Router(r *chi.Mux) *chi.Mux {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/livez", func(w http.ResponseWriter, r *http.Request) {
			a.writeJSON(w, http.StatusOK, map[string]string{"status": "OK"}, nil)
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Get("/scheduled", a.handleGetScheduledTasks)

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {

			})
		})
	})

	return r
}
