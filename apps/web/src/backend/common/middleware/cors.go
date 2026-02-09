package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/cors"
	"github.com/specvital/web/src/backend/common/config"
)

const (
	corsMaxAge     = 3600
	envCORSOrigins = "CORS_ORIGINS"
	envEnv         = "ENV"
	envProduction  = "production"
	wildcardOrigin = "*"
)

var (
	ErrCORSOriginsNotSet = errors.New("CORS_ORIGINS not set in production")
	ErrWildcardOrigin    = errors.New("wildcard origin not allowed with credentials")
)

func CORS(origins []string) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:   origins,
		ExposedHeaders:   []string{"Link"},
		MaxAge:           corsMaxAge,
	})
}

func GetAllowedOrigins() ([]string, error) {
	origins := os.Getenv(envCORSOrigins)
	if origins == "" {
		if os.Getenv(envEnv) == envProduction {
			return nil, ErrCORSOriginsNotSet
		}
		return []string{config.DefaultOrigin}, nil
	}

	var result []string
	for o := range strings.SplitSeq(origins, ",") {
		origin := strings.TrimSpace(o)
		if origin == wildcardOrigin {
			return nil, ErrWildcardOrigin
		}
		if origin != "" {
			result = append(result, origin)
		}
	}
	return result, nil
}
