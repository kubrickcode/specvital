package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCompress(t *testing.T) {
	smallBody := strings.Repeat("x", 500)  // Below threshold
	largeBody := strings.Repeat("x", 2000) // Above threshold

	tests := []struct {
		name           string
		acceptEncoding string
		responseBody   string
		expectGzip     bool
	}{
		{
			name:           "compress large response when gzip accepted",
			acceptEncoding: "gzip",
			responseBody:   largeBody,
			expectGzip:     true,
		},
		{
			name:           "no compress when gzip not accepted",
			acceptEncoding: "",
			responseBody:   largeBody,
			expectGzip:     false,
		},
		{
			name:           "no compress small response even with gzip accepted",
			acceptEncoding: "gzip",
			responseBody:   smallBody,
			expectGzip:     false,
		},
		{
			name:           "compress with gzip, deflate accepted",
			acceptEncoding: "gzip, deflate",
			responseBody:   largeBody,
			expectGzip:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.responseBody))
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", tt.acceptEncoding)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			gotGzip := rec.Header().Get("Content-Encoding") == "gzip"
			if gotGzip != tt.expectGzip {
				t.Errorf("expectGzip = %v, got %v", tt.expectGzip, gotGzip)
			}

			// Verify content is readable
			if gotGzip {
				gr, err := gzip.NewReader(rec.Body)
				if err != nil {
					t.Fatalf("failed to create gzip reader: %v", err)
				}
				defer gr.Close()

				decompressed, err := io.ReadAll(gr)
				if err != nil {
					t.Fatalf("failed to decompress: %v", err)
				}
				if string(decompressed) != tt.responseBody {
					t.Errorf("decompressed content mismatch")
				}
			} else {
				if rec.Body.String() != tt.responseBody {
					t.Errorf("body content mismatch")
				}
			}
		})
	}
}

func TestCompressPreservesStatusCode(t *testing.T) {
	largeBody := strings.Repeat("x", 2000)

	tests := []struct {
		name       string
		statusCode int
	}{
		{"ok", http.StatusOK},
		{"created", http.StatusCreated},
		{"bad request", http.StatusBadRequest},
		{"internal error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(largeBody))
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, rec.Code)
			}
		})
	}
}

func TestCompressVaryHeader(t *testing.T) {
	largeBody := strings.Repeat("x", 2000)

	handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(largeBody))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Vary") != "Accept-Encoding" {
		t.Errorf("expected Vary header to be 'Accept-Encoding', got %q", rec.Header().Get("Vary"))
	}
}

func TestCompressEmptyResponse(t *testing.T) {
	handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Empty response should not be gzipped
	if rec.Header().Get("Content-Encoding") == "gzip" {
		t.Error("empty response should not be gzipped")
	}
}

func TestShouldCompress(t *testing.T) {
	tests := []struct {
		name           string
		acceptEncoding string
		expected       bool
	}{
		{"gzip only", "gzip", true},
		{"gzip with deflate", "gzip, deflate", true},
		{"deflate only", "deflate", false},
		{"empty", "", false},
		{"br only", "br", false},
		{"all encodings", "gzip, deflate, br", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", tt.acceptEncoding)
			}

			got := shouldCompress(req)
			if got != tt.expected {
				t.Errorf("shouldCompress() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func BenchmarkCompress(b *testing.B) {
	largeBody := bytes.Repeat([]byte("benchmark test data "), 1000)

	handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(largeBody)
	}))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}
