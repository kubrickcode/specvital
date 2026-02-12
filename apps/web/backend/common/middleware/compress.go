package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

const minCompressLength = 1024

var gzipWriterPool = sync.Pool{
	New: func() any {
		return gzip.NewWriter(io.Discard)
	},
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf        *bytes.Buffer
	statusCode int
}

func (w *bufferedResponseWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w *bufferedResponseWriter) WriteHeader(code int) {
	w.statusCode = code
}

func Compress() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCompress(r) {
				next.ServeHTTP(w, r)
				return
			}

			buf := &bytes.Buffer{}
			bw := &bufferedResponseWriter{
				ResponseWriter: w,
				buf:            buf,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(bw, r)

			if buf.Len() < minCompressLength {
				w.WriteHeader(bw.statusCode)
				_, _ = w.Write(buf.Bytes())
				return
			}

			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")
			w.Header().Del("Content-Length")
			w.WriteHeader(bw.statusCode)

			gz := gzipWriterPool.Get().(*gzip.Writer)
			gz.Reset(w)

			_, writeErr := gz.Write(buf.Bytes())
			closeErr := gz.Close()

			if writeErr == nil && closeErr == nil {
				gzipWriterPool.Put(gz)
			} else {
				if writeErr != nil {
					slog.Warn("failed to write gzip data", "error", writeErr)
				}
				if closeErr != nil {
					slog.Warn("failed to close gzip writer", "error", closeErr)
				}
			}
		})
	}
}

func shouldCompress(r *http.Request) bool {
	acceptEncoding := r.Header.Get("Accept-Encoding")
	return strings.Contains(acceptEncoding, "gzip")
}
