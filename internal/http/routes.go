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
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				entries, err := a.resources.TaskResource.GetAllCrontabEntries()

				// Check what the error is at some point
				if err != nil {
					errRes := NewResponse(
						WithError[[]TaskDto](err),
					)

					a.writeJSON(w, http.StatusBadRequest, errRes, nil)
					return
				}

				var out []TaskDto = []TaskDto{}

				for _, item := range entries {
					i := TaskDto{
						ID:      item.ID,
						Command: item.Cmd,
						Cron:    item.Cron.String(),
					}

					out = append(out, i)
				}

				res := NewResponse(
					WithType[[]TaskDto]("task"),
					WithData(out),
				)

				a.writeJSON(w, http.StatusOK, res, nil)
			})

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {

			})
		})
	})

	return r
}
