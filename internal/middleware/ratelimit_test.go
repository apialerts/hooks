package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func setupRouter(rl *RateLimiter) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/{id}", func(r chi.Router) {
		r.Use(rl.Limit)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})
	return r
}

func TestRateLimiterAllowsRequests(t *testing.T) {
	rl := NewRateLimiter()
	router := setupRouter(rl)

	req := httptest.NewRequest("GET", "/test-endpoint/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRateLimiterBlocksAfterLimit(t *testing.T) {
	rl := NewRateLimiter()
	router := setupRouter(rl)

	for i := 0; i < 60; i++ {
		req := httptest.NewRequest("GET", "/test-endpoint/", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("request %d should be allowed, got %d", i+1, w.Code)
		}
	}

	// 61st request should be blocked
	req := httptest.NewRequest("GET", "/test-endpoint/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
}

func TestRateLimiterSeparatesEndpoints(t *testing.T) {
	rl := NewRateLimiter()
	router := setupRouter(rl)

	// Exhaust limit on endpoint A
	for i := 0; i < 61; i++ {
		req := httptest.NewRequest("GET", "/endpoint-a/", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Endpoint B should still work
	req := httptest.NewRequest("GET", "/endpoint-b/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("endpoint-b should not be rate limited, got %d", w.Code)
	}
}

func TestRateLimiterSeparatesIPs(t *testing.T) {
	rl := NewRateLimiter()
	router := setupRouter(rl)

	// Exhaust limit for IP A
	for i := 0; i < 61; i++ {
		req := httptest.NewRequest("GET", "/test-endpoint/", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Different IP should still work
	req := httptest.NewRequest("GET", "/test-endpoint/", nil)
	req.RemoteAddr = "10.0.0.2:5678"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("different IP should not be rate limited, got %d", w.Code)
	}
}
