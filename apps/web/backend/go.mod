module github.com/kubrickcode/specvital/apps/web/backend

go 1.24.0

require (
	github.com/bradleyfalzon/ghinstallation/v2 v2.17.0
	github.com/cockroachdb/errors v1.11.3
	github.com/go-chi/chi/v5 v5.2.0
	github.com/go-chi/cors v1.2.1
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/go-github/v75 v75.0.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/joho/godotenv v1.5.1
	github.com/kubrickcode/specvital/lib v0.0.0
	github.com/oapi-codegen/runtime v1.1.2
	github.com/riverqueue/river v0.26.0
	github.com/riverqueue/river/riverdriver/riverpgxv5 v0.26.0
	github.com/riverqueue/river/rivertype v0.28.0
	golang.org/x/oauth2 v0.34.0
	golang.org/x/sync v0.19.0
	golang.org/x/time v0.14.0
)

replace github.com/kubrickcode/specvital/lib => ../../../lib

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/getsentry/sentry-go v0.27.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/riverqueue/river/riverdriver v0.28.0 // indirect
	github.com/riverqueue/river/rivershared v0.26.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	go.uber.org/goleak v1.3.0 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
