package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type RateLimiter struct {
	mu       sync.Mutex
	counters map[string]*counter
}

type counter struct {
	count   int
	resetAt time.Time
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		counters: make(map[string]*counter),
	}

	// Cleanup old entries every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.mu.Lock()
			now := time.Now()
			for k, c := range rl.counters {
				if now.After(c.resetAt) {
					delete(rl.counters, k)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		key := id + ":" + r.RemoteAddr

		rl.mu.Lock()
		c, exists := rl.counters[key]
		now := time.Now()

		if !exists || now.After(c.resetAt) {
			c = &counter{
				count:   0,
				resetAt: now.Add(1 * time.Minute),
			}
			rl.counters[key] = c
		}

		c.count++
		count := c.count
		rl.mu.Unlock()

		if count > 60 {
			http.Error(w, "rate limit exceeded (60 requests/minute)", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
