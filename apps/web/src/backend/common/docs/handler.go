package docs

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/specvital/web/src/backend/api"
)

//go:embed index.html
var indexHTML []byte

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/docs", h.serveIndex)
	r.Get("/api/openapi.yaml", h.serveSpec)
}

func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML)
}

func (h *Handler) serveSpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.Write(api.OpenAPISpec)
}
