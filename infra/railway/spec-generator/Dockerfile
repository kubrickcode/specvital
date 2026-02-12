# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY apps/worker/go.mod apps/worker/go.sum ./
COPY lib/go.mod lib/go.sum /lib/

RUN go mod download

COPY apps/worker/ ./
COPY lib/ /lib/

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /service ./cmd/spec-generator

FROM alpine:3.21

RUN apk add --no-cache ca-certificates git

RUN adduser -D -u 1000 appuser

WORKDIR /app

COPY --from=builder /service .

USER appuser

ENTRYPOINT ["./service"]
