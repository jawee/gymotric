package ratelimiter

import (
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	// Mutex to protect concurrent access to the requests map
	mu sync.Mutex
	// Map to store request counts and window start time per IP
	requests map[string]*clientStats
	// Duration of the fixed window
	window time.Duration
	// Maximum number of requests allowed per window
	limit int
}
type clientStats struct {
	// Time when the current window started
	windowStart time.Time
	// Number of requests made within the current window
	count int
}

func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]*clientStats),
		window:   window,
		limit:    limit,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	stats, ok := rl.requests[ip]

	// If no stats for this IP, create new stats
	if !ok || time.Since(stats.windowStart) >= rl.window {
		rl.requests[ip] = &clientStats{
			windowStart: time.Now(),
			count:       1,
		}
		return true
	}

	// If within the window and count is less than the limit, increment count and allow
	if stats.count < rl.limit {
		stats.count++
		return true
	}

	// If within the window and count is at or above the limit, do not allow
	return false
}

func RateLimitMiddleware(limiter *RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Real-IP")
		slog.Debug("Rate limit middleware invoked", "path", r.URL.Path, "IP", ip)

		if !limiter.Allow(ip) {
			slog.Warn("Rate limit exceeded for IP", "ip", ip, "path", r.URL.Path)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Request is allowed, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

