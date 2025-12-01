package coco_http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *app) apiV1Router(r *chi.Mux) *chi.Mux {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/livez", func(w http.ResponseWriter, r *http.Request) {
			res := NewGetTasksResponse([]TaskDto{})

			fmt.Fprintf(w, "hello world with %s", res.Type)
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				res := NewGetTasksResponse([]TaskDto{})

				fmt.Fprintf(w, "hello world with %s", res.Type)
			})

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {

			})
		})
	})

	return r
}
