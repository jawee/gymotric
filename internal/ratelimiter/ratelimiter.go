package ratelimiter

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]*clientStats
	window   time.Duration
	limit    int
}
type clientStats struct {
	windowStart time.Time
	count       int
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

	if !ok || time.Since(stats.windowStart) >= rl.window {
		rl.requests[ip] = &clientStats{
			windowStart: time.Now(),
			count:       1,
		}
		return true
	}

	if stats.count < rl.limit {
		stats.count++
		return true
	}

	return false
}

func (rl *RateLimiter) GetWindowExpiration() time.Duration {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if len(rl.requests) == 0 {
		return 0
	}

	var earliestExpiration time.Time
	for _, stats := range rl.requests {
		if earliestExpiration.IsZero() || stats.windowStart.Add(rl.window).Before(earliestExpiration) {
			earliestExpiration = stats.windowStart.Add(rl.window)
		}
	}

	return time.Until(earliestExpiration)
}

func RateLimitMiddleware(limiter *RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Real-IP")
		slog.Debug("Rate limit middleware invoked", "path", r.URL.Path, "IP", ip)

		if ip == "" {
			slog.Warn("No IP address found in request header", "path", r.URL.Path)
			remoteAddr := r.RemoteAddr
			host, _, err := net.SplitHostPort(remoteAddr)
			if err != nil {
				slog.Error("Failed to parse remote address", "remoteAddr", remoteAddr, "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			ip = host
		}

		if !limiter.Allow(ip) {
			stats, ok := limiter.requests[ip]
			if ok {
				timeUntilReset := max(stats.windowStart.Add(limiter.window).Sub(time.Now()), 0)
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(timeUntilReset.Seconds())))
			}
			slog.Warn("Rate limit exceeded for IP", "ip", ip, "path", r.URL.Path)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
