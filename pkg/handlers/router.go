package handlers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(logger *logrus.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(NewStructuredLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(PrometheusHandler())
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "sup")
	})

	i := chi.NewRouter()
	i.Get("/internal/healthz", noContent())
	i.Get("/internal/metrics", promhttp.Handler().ServeHTTP)

	r.Mount("/internal", i)

	return r
}

func noContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.NoContent(w, r)
	}
}
