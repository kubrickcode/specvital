# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY apps/web/src/backend/go.mod apps/web/src/backend/go.sum ./

RUN go mod download

COPY apps/web/src/backend/ .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates git

RUN adduser -D -u 1000 appuser

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8000

USER appuser

ENTRYPOINT ["/app/server"]
