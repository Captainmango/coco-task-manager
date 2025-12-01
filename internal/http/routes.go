package coco_http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a *app) apiV1Router(r *chi.Mux) *chi.Mux {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/livez", func(w http.ResponseWriter, r *http.Request) {
			res := NewResponse(
				WithType[[]TaskDto]("task"),
				WithData([]TaskDto{
					{
						ID: uuid.NullUUID{}.UUID,
						Command: "test",
					},
				}),
			)

			a.writeJSON(w, http.StatusOK, res, nil)
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {

			})

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {

			})
		})
	})

	return r
}
