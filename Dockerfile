# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY src/go.mod src/go.sum ./

RUN go mod download

COPY src/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /collector ./cmd/collector

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

RUN adduser -D -u 1000 collector

WORKDIR /app

COPY --from=builder /collector .

USER collector

ENTRYPOINT ["./collector"]
