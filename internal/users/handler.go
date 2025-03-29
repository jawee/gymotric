package users

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
	"weight-tracker/internal/database"
	"weight-tracker/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

type createUserAndReturnIdRequest struct {
	Username string
	Password string
}

type loginRequest struct {
	Username string
	Password string
}

type changePasswordRequest struct {
	OldPassword string
	NewPassword string
}

type handler struct {
	service Service
}

func AddEndpoints(mux *http.ServeMux, s database.Service, authenticationWrapper func(next http.Handler) http.Handler) {
	handler := handler{
		service: NewService(&usersRepository{s.GetRepository()}),
	}

	mux.Handle("POST /users", http.HandlerFunc(handler.createUserHandler))
 
	mux.Handle("POST /auth/login", http.HandlerFunc(handler.loginHandler))
	mux.Handle("POST /auth/token", http.HandlerFunc(handler.refreshHandler))

	mux.Handle("GET /me", authenticationWrapper(http.HandlerFunc(handler.meHandler)))
	mux.Handle("PUT /me/password", authenticationWrapper(http.HandlerFunc(handler.changePasswordHandler)))

	mux.Handle("POST /logout", authenticationWrapper(http.HandlerFunc(handler.logoutHandler)))
}

func (s *handler) changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)

	decoder := json.NewDecoder(r.Body)
	var request changePasswordRequest
	err := decoder.Decode(&request)

	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = s.service.ChangePassword(r.Context(), request, userId)
	if err != nil {
		slog.Error("Failed to change password", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	cookie := createCookie(utils.AccessTokenCookieName, "", time.Now().Add(time.Second))
	refresh_cookie := createCookie(utils.RefreshTokenCookieName, "", time.Now().Add(time.Second))

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &refresh_cookie)

	w.WriteHeader(http.StatusNoContent)
}

func (s *handler) meHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)

	user, err := s.service.GetByUserId(r.Context(), userId)
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(user)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
	}
}

func (s *handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := createCookie(utils.AccessTokenCookieName, "", time.Now().Add(time.Second))
	refresh_cookie := createCookie(utils.RefreshTokenCookieName, "", time.Now().Add(time.Second))

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &refresh_cookie)
	w.Header().Set("Content-Type", "application/json")
}

func createRefreshToken(userId string) (string, error) {
	signingKey := os.Getenv(utils.EnvJwtRefreshSignKey)
	tokenExpiration, err := strconv.Atoi(os.Getenv(utils.EnvJwtRefreshExpireMinutes))
	if err != nil {
		slog.Error("Failed to convert JWT_EXPIRATION to int", "error", err)
		return "", err
	}

	mySigningKey := []byte(signingKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * time.Duration(tokenExpiration))),
		Issuer:    "weight-tracker",
		Subject:   userId,
		Audience:  []string{"weight-tracker"},
	})
	return token.SignedString(mySigningKey)
}

func createCookie(name string, value string, expiration time.Time) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  expiration,
		SameSite: http.SameSiteLaxMode,
	}
}

func getSubjectFromCookie(cookieName string, signingKey string, r *http.Request) (string, error) {
	cookieTokenStr := ""
	for _, cookie := range r.Cookies() {
		if cookie.Name == cookieName {
			cookieTokenStr = cookie.Value
		}
	}

	if cookieTokenStr != "" {
		cookieToken, err := jwt.Parse(cookieTokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		})

		if err != nil {
			slog.Error("Cookie: refresh token error", "error", err)
			return "", err
		}

		if claims, ok := cookieToken.Claims.(jwt.MapClaims); ok {
			sub, err := claims.GetSubject()
			if err != nil {
				slog.Error("Cookie: GetSubject", "error", err)
				return "", err
			}
			return sub, nil
		}
		slog.Error("Cookie: error getting claims", "error", err)
		return "", errors.New("error getting claims")
	}

	return "", errors.New("No token found")
}

func (s *handler) createTokenResponse(w http.ResponseWriter, sub string) error {
	tokenExpiration, err := strconv.Atoi(os.Getenv(utils.EnvJwtExpireMinutes))


	if err != nil {
		slog.Error("Failed to convert JWT_EXPIRE_MINUTES to int", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return err
	}

	_, err = s.service.GetByUserId(context.Background(), sub)
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return err
	}

	newToken, err := s.service.CreateToken(sub)

	if err != nil {
		slog.Error("Failed to create new token", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return err
	}

	cookie := createCookie(utils.AccessTokenCookieName, newToken, time.Now().Add(time.Minute*time.Duration(tokenExpiration)))

	refresh_token, err := createRefreshToken(sub)
	if err != nil {
		slog.Error("Failed to create refresh token", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return err
	}

	refresh_cookie := createCookie(utils.RefreshTokenCookieName, refresh_token, time.Now().Add(time.Hour*24))

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &refresh_cookie)

	resp := map[string]interface{}{
		"access_token":  newToken,
		"token_type":    "Bearer",
		"expires_in":    tokenExpiration * 60,
		"refresh_token": refresh_token,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
		return err
	}

	return nil
}

func (s *handler) refreshHandler(w http.ResponseWriter, r *http.Request) {
	signingKey := os.Getenv(utils.EnvJwtRefreshSignKey)

	cookieSub, err := getSubjectFromCookie(utils.RefreshTokenCookieName, signingKey, r)
	if err == nil {
		err = s.createTokenResponse(w, cookieSub)

		if err != nil {
			slog.Error("request failed authentication", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		return
	}

	slog.Error("Cookie: Failed to get subject from refresh token", "error", err)
	slog.Info("Cookie failed. Trying query parameter token")

	//body
	refreshToken := r.URL.Query().Get("refresh_token")
	if refreshToken == "" {
		slog.Error("No refresh token found")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
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
			slog.Error("Couldn't get sub claim from token", "error", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		slog.Info("Success", "sub", sub)

		err = s.createTokenResponse(w, sub)

		if err != nil {
			slog.Error("request failed authentication", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		return
	}
	slog.Error("error getting claims", "error", err)
	w.WriteHeader(http.StatusUnauthorized)
	return
}

func (s *handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	tokenExpiration, err := strconv.Atoi(os.Getenv(utils.EnvJwtExpireMinutes))
	if err != nil {
		slog.Error("Failed to convert JWT_EXPIRE_MINUTES to int", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var t loginRequest
	err = decoder.Decode(&t)

	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := s.service.Login(r.Context(), t)

	if err != nil {
		slog.Warn("Failed to login", "error", err)
		http.Error(w, "Failed to login", http.StatusBadRequest)
		return
	}

	cookie := createCookie(utils.AccessTokenCookieName, token.Token, time.Now().Add(time.Minute*time.Duration(tokenExpiration)))

	refresh_token, err := createRefreshToken(token.UserId)
	if err != nil {
		slog.Warn("Failed to create refresh token", "error", err)
		http.Error(w, "Failed to login", http.StatusBadRequest)
		return
	}

	refresh_cookie := createCookie(utils.RefreshTokenCookieName, refresh_token, time.Now().Add(time.Hour*24))

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &refresh_cookie)

	resp := map[string]interface{}{
		"access_token":  token,
		"token_type":    "Bearer",
		"expires_in":    tokenExpiration * 60,
		"refresh_token": "refresh_token",
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
	}
}

func (s *handler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t createUserAndReturnIdRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Warn("Invalid request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := s.service.CreateAndReturnId(r.Context(), t)

	if err != nil {
		slog.Warn("Failed to create user", "error", err)
		http.Error(w, "Failed to create user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	resp := map[string]interface{}{"id": id}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
	}
}
