package http

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	v1 "github.com/swift1337/webmajor/internal/server/http/v1"
)

type Opt func(r chi.Router)

func NewRouter(opts ...Opt) chi.Router {
	r := chi.NewRouter()

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func WithDashboardAssets(path string, files http.FileSystem) Opt {
	return func(r chi.Router) {
		r.Handle(path+"/*", http.StripPrefix(path, http.FileServer(files)))
	}
}

func WithDashboardAPI(path string, handler *v1.Handler) Opt {
	return func(r chi.Router) {
		api := chi.NewRouter()
		api.Get("/request", handler.ListRequests)

		r.Mount(path, api)
	}
}

func WithProxy(handler *v1.Handler) Opt {
	return func(r chi.Router) {
		r.Handle("/*", handler)
	}
}
