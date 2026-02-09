package server

import "github.com/go-chi/chi/v5"

type RouteRegistrar interface {
	RegisterRoutes(r chi.Router)
}
