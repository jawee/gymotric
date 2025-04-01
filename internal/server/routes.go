package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"weight-tracker/internal/exercises"
	"weight-tracker/internal/exercisetypes"
	"weight-tracker/internal/sets"
	"weight-tracker/internal/statistics"
	"weight-tracker/internal/users"
	"weight-tracker/internal/utils"
	"weight-tracker/internal/workouts"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	_ "github.com/joho/godotenv/autoload"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /health", http.HandlerFunc(s.healthHandler))

	users.AddEndpoints(mux, s.db, s.AuthenticatedMiddleware)

	exercisetypes.AddEndpoints(mux, s.db, s.AuthenticatedMiddleware)

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
			cookieToken, err := jwt.Parse(cookieTokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(signingKey), nil
			})

			if err != nil {
				slog.Error("Cookie: request failed authentication", "error", err)
			}

			if err == nil {
				if claims, ok := cookieToken.Claims.(jwt.MapClaims); ok {
					sub, err := claims.GetSubject()
					if err != nil {
						slog.Error("Cookie: GetSubject", "error", err)
						return
					}
					claimsCtx := context.WithValue(r.Context(), "sub", sub)
					r = r.WithContext(claimsCtx)
					slog.Info("Cookie: Success", "sub", sub)
					next.ServeHTTP(w, r)
					return
				} else {
					slog.Error("Cookie: error getting claims", "error", err)
				}
			}
		}

		slog.Info("Cookie failed. Trying header token")

		//header token
		extractor := request.AuthorizationHeaderExtractor
		token, err := request.ParseFromRequest(r, extractor, func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		})

		if err != nil {
			slog.Error("request failed authentication", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			sub, err := claims.GetSubject()
			if err != nil {
				slog.Error("GetSubject", "error", err)
				return
			}
			claimsCtx := context.WithValue(r.Context(), "sub", sub)
			r = r.WithContext(claimsCtx)
			slog.Info("Success", "sub", sub)
		} else {
			slog.Error("error getting claims", "error", err)
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request", "Path", r.URL.Path, "Method", r.Method)

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
			slog.Info("Returning 204")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		slog.Info("Proceeding")
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
