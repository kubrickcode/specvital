package logger

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

type Logger struct {
	base  *slog.Logger
	attrs []any
}

func New() *Logger {
	return &Logger{base: slog.Default()}
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		base:  l.base,
		attrs: append(l.attrs, args...),
	}
}

func (l *Logger) logger(ctx context.Context) *slog.Logger {
	allAttrs := make([]any, 0, len(l.attrs)+2)
	allAttrs = append(allAttrs, "request_id", middleware.GetReqID(ctx))
	allAttrs = append(allAttrs, l.attrs...)
	return l.base.With(allAttrs...)
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger(ctx).Debug(msg, args...)
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.logger(ctx).Info(msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger(ctx).Warn(msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.logger(ctx).Error(msg, args...)
}
