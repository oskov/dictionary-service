package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/oskov/dictionary-service/internal/api/http/oapi"
	"github.com/oskov/dictionary-service/internal/application"
)

func NewServer(ctx context.Context, port uint, app application.App) *http.Server {
	api := NewAPI(app)
	r := createRouter(api)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second * 30,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	return srv
}

func createRouter(api *API) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Group(func(r chi.Router) {
		strictHandler := oapi.NewStrictHandlerWithOptions(api, nil, oapi.StrictHTTPServerOptions{})

		oapi.HandlerFromMux(strictHandler, r)
	})

	return r
}
