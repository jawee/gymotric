package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
	"weight-tracker/internal/exercises"
	"weight-tracker/internal/exerciseitems"
	"weight-tracker/internal/exercisetypes"
	"weight-tracker/internal/ratelimiter"
	"weight-tracker/internal/repository"
	"weight-tracker/internal/sets"
	"weight-tracker/internal/statistics"
	"weight-tracker/internal/users"
	"weight-tracker/internal/utils"
	"weight-tracker/internal/workouts"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

func (s *Server) RegisterRoutes() http.Handler {
	rateLimiter := ratelimiter.NewRateLimiter(1*time.Minute, 500)
	mux := http.NewServeMux()

	mux.Handle("GET /health", http.HandlerFunc(s.healthHandler))
	mux.Handle("GET /ip", http.HandlerFunc(s.ipHandler))

	users.AddEndpoints(mux, s.db, s.AuthenticatedMiddleware, ratelimiter.RateLimitMiddleware, rateLimiter)

	exercisetypes.AddEndpoints(mux, s.db, s.AuthenticatedMiddleware)

	exerciseitems.AddEndpoints(mux, s.db, s.AuthenticatedMiddleware)

	workouts.AddEndpoints(mux, s.db, s.AuthenticatedMiddleware)

	sets.AddEndpoints(mux, s.db, s.AuthenticatedMiddleware)

	exercises.AddEndpoints(mux, s.db, s.AuthenticatedMiddleware)

	statistics.AddEndpoints(mux, s.db, s.AuthenticatedMiddleware)

	return s.corsMiddleware(s.loggingMiddleware(mux))
}

func (s *Server) AuthenticatedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signingKey := os.Getenv(utils.EnvJwtSignKey)

		cookieTokenStr := ""
		for _, cookie := range r.Cookies() {
			if cookie.Name == utils.AccessTokenCookieName {
				cookieTokenStr = cookie.Value
			}
		}

		if cookieTokenStr != "" {
			args := repository.CheckIfTokenExistsParams{
				Token:     cookieTokenStr,
				TokenType: "access_token",
			}
			slog.Debug("Cookie: CheckIfTokenExists", "token", cookieTokenStr)
			id, err := s.db.GetRepository().CheckIfTokenExists(r.Context(), args)
			if err != nil {
				if err == sql.ErrNoRows {
					slog.Debug("Cookie: Token not found in expired tokens", "token", cookieTokenStr)
					// Token not found, proceed with authentication
				} else {
					slog.Error("Cookie: CheckIfTokenExists failed", "error", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			if id > 0 {
				slog.Info("Cookie: Token is expired", "token", cookieTokenStr)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			refreshTokenStr := ""
			for _, cookie := range r.Cookies() {
				if cookie.Name == utils.RefreshTokenCookieName {
					refreshTokenStr = cookie.Value
				}
			}

			cookieToken, err := jwt.Parse(cookieTokenStr, func(token *jwt.Token) (any, error) {
				return []byte(signingKey), nil
			})

			if err != nil {
				slog.Error("Cookie: request failed authentication", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if claims, ok := cookieToken.Claims.(jwt.MapClaims); ok {
				sub, err := claims.GetSubject()
				if err != nil {
					slog.Error("Cookie: GetSubject", "error", err)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				claimsCtx := context.WithValue(r.Context(), "sub", sub)
				claimsCtx = context.WithValue(claimsCtx, "access_token", cookieTokenStr)
				claimsCtx = context.WithValue(claimsCtx, "refresh_token", refreshTokenStr)
				r = r.WithContext(claimsCtx)
				slog.Debug("Cookie: Success", "sub", sub)
				next.ServeHTTP(w, r)
				return
			}

			slog.Error("Cookie: error getting claims", "error", err)
		}

		w.WriteHeader(http.StatusUnauthorized)
		return
	})
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Request", "Path", r.URL.Path, "Method", r.Method)

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, ApiKey")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			slog.Debug("Returning 204")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		slog.Debug("Proceeding")
		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-wt-api-key")
	if apiKey != os.Getenv(utils.EnvApiKey) {
		slog.Error("Invalid API key", "provided", apiKey)
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		slog.Error("Failed to marshal health check response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		slog.Warn("Failed to write response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
	}
}

func (s *Server) ipHandler(w http.ResponseWriter, r *http.Request) {
	ipFromHeader := r.Header.Get("X-Real-IP")

	if ipFromHeader != "" {
		w.Write([]byte(ipFromHeader))
		return
	}
	slog.Error("X-Real-IP header not found. Trying remoteAddr")

	ipFromRemoteAddr, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		slog.Error("Failed to split remote address", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(ipFromRemoteAddr))
}
