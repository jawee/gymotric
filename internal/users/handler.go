package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
	"weight-tracker/internal/database"
	"weight-tracker/internal/email"
	"weight-tracker/internal/ratelimiter"
	"weight-tracker/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

type registrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

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

type resetPasswordRequest struct {
	Email string
}

type resetPasswordConfirmRequest struct {
	Token    string
	Password string
}

type changeEmailRequest struct {
	Email string
}

type handler struct {
	service Service
}

func AddEndpoints(
	mux *http.ServeMux,
	s database.Service,
	authenticationWrapper func(next http.Handler) http.Handler,
	rateLimitWrapper func(limiter *ratelimiter.RateLimiter, next http.Handler) http.Handler,
	rateLimiter *ratelimiter.RateLimiter,
) {

	handler := handler{
		service: NewService(&usersRepository{s.GetRepository()}),
	}

	mux.Handle("POST /users", authenticationWrapper(http.HandlerFunc(handler.createUserHandler)))

	mux.Handle("POST /auth/login", rateLimitWrapper(rateLimiter, http.HandlerFunc(handler.loginHandler)))
	mux.Handle("POST /auth/token", rateLimitWrapper(rateLimiter, http.HandlerFunc(handler.refreshHandler)))

	mux.Handle("GET /me", authenticationWrapper(http.HandlerFunc(handler.meHandler)))
	mux.Handle("PUT /me/password", authenticationWrapper(http.HandlerFunc(handler.changePasswordHandler)))
	mux.Handle("PUT /me/email", authenticationWrapper(http.HandlerFunc(handler.changeEmailHandler)))

	mux.Handle("POST /confirm-email", rateLimitWrapper(rateLimiter, http.HandlerFunc(handler.confirmEmailHandler)))

	mux.Handle("POST /reset-password", rateLimitWrapper(rateLimiter, http.HandlerFunc(handler.resetPasswordHandler)))
	mux.Handle("POST /reset-password/confirm", rateLimitWrapper(rateLimiter, http.HandlerFunc(handler.resetPasswordConfirmHandler)))

	mux.Handle("POST /logout", authenticationWrapper(http.HandlerFunc(handler.logoutHandler)))

	mux.Handle("POST /register", rateLimitWrapper(rateLimiter, http.HandlerFunc(handler.registrationHandler)))
	mux.Handle("POST /register/confirm", rateLimitWrapper(rateLimiter, http.HandlerFunc(handler.confirmRegistrationHandler)))
}

func (s *handler) registrationHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request registrationRequest
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	id, err := s.service.Register(r.Context(), request)

	if err != nil {
		slog.Error("Failed to register", "error", err, "email", request.Email)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token, err := s.service.CreateAccountConfirmationToken(r.Context(), id)
	if err != nil {
		slog.Error("Failed to create confirmation token", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	baseUrl := os.Getenv("BASE_URL")
	err = email.SendAccountConfirmation(request.Email, email.SendAccountConfirmationData{
		Name: request.Username,
		Link: baseUrl + "/confirm-registration?token=" + token,
	})

	w.WriteHeader(http.StatusCreated)
}

func (s *handler) confirmRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Confirm registration handler")
	tokenString := r.URL.Query().Get("token")

	if tokenString == "" {
		slog.Error("No token found")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		signingKey := os.Getenv(utils.EnvJwtSignKey)
		return []byte(signingKey), nil
	})

	if err != nil {
		slog.Error("Failed to parse token", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	slog.Debug("Parsed token")
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		sub, err := claims.GetSubject()

		if err != nil {
			slog.Error("GetSubject", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		slog.Debug("Success", "sub", sub)
		err = s.service.ConfirmAccount(r.Context(), sub)
		if err != nil {
			slog.Error("Failed to confirm email", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *handler) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request resetPasswordRequest
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	user, err := s.service.GetByEmail(r.Context(), request.Email)
	if err != nil {
		slog.Error("Failed to get user", "error", err, "email", request.Email)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token, err := s.service.CreateResetPasswordToken(r.Context(), user.ID)
	if err != nil {
		slog.Error("Failed to create reset password token", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	baseUrl := os.Getenv("BASE_URL")

	err = email.SendPasswordReset(request.Email, email.ResetPasswordEmailData{
		Name:      user.Username,
		ResetLink: baseUrl + "/password-reset/" + token,
	})

	if err != nil {
		slog.Error("Failed to send email", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *handler) resetPasswordConfirmHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Reset password confirm handler")
	decoder := json.NewDecoder(r.Body)
	var request resetPasswordConfirmRequest
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if request.Password == "" {
		slog.Error("Password is empty")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if request.Token == "" {
		slog.Error("Token is empty")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (any, error) {
		signingKey := os.Getenv(utils.EnvJwtSignKey)
		return []byte(signingKey), nil
	})

	if err != nil {
		slog.Error("Failed to parse token", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	slog.Debug("Parsed token")
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		sub, err := claims.GetSubject()

		if err != nil {
			slog.Error("GetSubject", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		slog.Debug("Success", "sub", sub)
		err = s.service.ResetPassword(r.Context(), sub, request.Password)
		if err != nil {
			slog.Error("Failed to reset password", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (s *handler) confirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Confirm email handler")
	tokenString := r.URL.Query().Get("token")

	if tokenString == "" {
		slog.Error("No token found")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &emailConfirmationCustomClaims{}, func(token *jwt.Token) (any, error) {
		signingKey := os.Getenv(utils.EnvJwtSignKey)
		return []byte(signingKey), nil
	})

	if err != nil {
		slog.Error("Failed to parse token", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	slog.Debug("Parsed token")
	if claims, ok := token.Claims.(*emailConfirmationCustomClaims); ok {
		sub, err := claims.GetSubject()

		if err != nil {
			slog.Error("GetSubject", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		email, err := claims.GetEmail()
		if err != nil {
			slog.Error("GetEmail", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		slog.Debug("Success", "sub", sub, "email", email)
		err = s.service.ConfirmEmail(r.Context(), sub, email)
		if err != nil {
			slog.Error("Failed to confirm email", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *handler) changeEmailHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	decoder := json.NewDecoder(r.Body)
	var request changeEmailRequest
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	user, err := s.service.GetByUserId(r.Context(), userId)
	if err != nil {
		slog.Error("Failed to get user", "error", err, "userId", userId)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if user.Email == request.Email {
		slog.Error("Email is the same as current email")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token, err := s.service.CreateConfirmationToken(r.Context(), userId, request.Email)
	if err != nil {
		slog.Error("Failed to create confirmation token", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	baseUrl := os.Getenv("BASE_URL")

	err = email.SendEmailConfirmation(request.Email, email.SendEmailConfirmationData{
		Name: user.Username,
		Link: baseUrl + "/confirm-email?token=" + token,
	})

	if err != nil {
		slog.Error("Failed to send email", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

	utils.ReturnJson(w, jsonResp)
}

func (s *handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	userToken := r.Context().Value("access_token").(string)
	refreshToken := r.Context().Value("refresh_token").(string)

	err := s.service.Logout(r.Context(), userToken, refreshToken) 
	if err != nil {
		slog.Error("Failed to logout", "error", err)
	}

	cookie := createCookie(utils.AccessTokenCookieName, "", time.Now())
	refresh_cookie := createCookie(utils.RefreshTokenCookieName, "", time.Now())

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &refresh_cookie)
	w.Header().Set("Content-Type", "application/json")
}

func createRefreshToken(userId string) (string, error) {
	signingKey := os.Getenv(utils.EnvJwtRefreshSignKey)
	tokenExpiration, err := strconv.Atoi(os.Getenv(utils.EnvJwtRefreshExpireMinutes))
	if err != nil {
		slog.Error("Failed to convert JWT_EXPIRATION to int", "error", err)
		return "", fmt.Errorf("Failed to convert JWT_EXPIRATION to int: %w", err)
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

// TODO: currently only used for refresh token, investigate.
func (s *handler) getSubjectFromCookie(ctx context.Context, cookieName, signingKey, tokenType string, cookies []*http.Cookie) (string, error) {
	cookieTokenStr := ""
	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			cookieTokenStr = cookie.Value
		}
	}

	if cookieTokenStr != "" {
		valid := s.service.IsTokenValid(ctx, cookieTokenStr, tokenType) 
		if !valid {
			slog.Error("Cookie: Token is invalid", "cookieName", cookieName, "token", cookieTokenStr)
			return "", errors.New("token is invalid")
		}
		cookieToken, err := jwt.Parse(cookieTokenStr, func(token *jwt.Token) (any, error) {
			return []byte(signingKey), nil
		})

		if err != nil {
			slog.Error("Cookie: refresh token error", "error", err)
			return "", fmt.Errorf("Cookie: refresh token error: %w", err)
		}

		if claims, ok := cookieToken.Claims.(jwt.MapClaims); ok {
			sub, err := claims.GetSubject()
			if err != nil {
				slog.Error("Cookie: GetSubject", "error", err)
				return "", fmt.Errorf("Cookie: GetSubject error: %w", err)
			}
			return sub, nil
		}
		slog.Error("Cookie: error getting claims", "error", err)
		return "", errors.New("error getting claims")
	}

	return "", NoTokenFoundError
}

var NoTokenFoundError = errors.New("No token found")

func (s *handler) createTokenResponse(w http.ResponseWriter, sub string) error {
	tokenExpiration, err := strconv.Atoi(os.Getenv(utils.EnvJwtExpireMinutes))

	if err != nil {
		slog.Error("Failed to convert JWT_EXPIRE_MINUTES to int", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return fmt.Errorf("Failed to convert JWT_EXPIRE_MINUTES to int: %w", err)
	}

	_, err = s.service.GetByUserId(context.Background(), sub)
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return fmt.Errorf("Failed to get user: %w", err)
	}

	newToken, err := s.service.CreateToken(sub)

	if err != nil {
		slog.Error("Failed to create new token", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return fmt.Errorf("Failed to create new token: %w", err)
	}

	cookie := createCookie(utils.AccessTokenCookieName, newToken, time.Now().Add(time.Minute*time.Duration(tokenExpiration)))

	refresh_token, err := createRefreshToken(sub)
	if err != nil {
		slog.Error("Failed to create refresh token", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return fmt.Errorf("Failed to create refresh token: %w", err)
	}

	refresh_cookie := createCookie(utils.RefreshTokenCookieName, refresh_token, time.Now().Add(time.Hour*24))

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &refresh_cookie)
	return nil
}

func (s *handler) refreshHandler(w http.ResponseWriter, r *http.Request) {
	signingKey := os.Getenv(utils.EnvJwtRefreshSignKey)

	cookieSub, err := s.getSubjectFromCookie(r.Context(), utils.RefreshTokenCookieName, signingKey, "refresh_token", r.Cookies())
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
	slog.Error("error getting claims", "error", err)
	w.WriteHeader(http.StatusUnauthorized)
	return
}

func (s *handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	tokenExpiration, err := strconv.Atoi(os.Getenv(utils.EnvJwtExpireMinutes))
	slog.Debug("Login handler invoked")
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

	loginResponse, err := s.service.Login(r.Context(), t)

	if err != nil {
		slog.Warn("Failed to login", "error", err)
		http.Error(w, "Failed to login", http.StatusBadRequest)
		return
	}

	cookie := createCookie(utils.AccessTokenCookieName, loginResponse.Token, time.Now().Add(time.Minute*time.Duration(tokenExpiration)))

	refreshToken, err := createRefreshToken(loginResponse.UserId)
	if err != nil {
		slog.Warn("Failed to create refresh token", "error", err)
		http.Error(w, "Failed to login", http.StatusBadRequest)
		return
	}

	refresh_cookie := createCookie(utils.RefreshTokenCookieName, refreshToken, time.Now().Add(time.Hour*24))

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &refresh_cookie)
	w.WriteHeader(http.StatusOK)
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

	jsonResp, err := utils.CreateIdResponse(id)

	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	utils.ReturnJson(w, jsonResp)
}
