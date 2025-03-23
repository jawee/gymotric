package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"weight-tracker/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

func createToken(t *testing.T, userId string, expiration int) string {
	signingKey := "sekrit"

	mySigningKey := []byte(signingKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * time.Duration(expiration))),
		Issuer:    "weight-tracker",
		Subject:   userId,
		Audience:  []string{"weight-tracker"},
	})
	tokStr, err := token.SignedString(mySigningKey)
	if err != nil {
		t.Fatal("Failed to sign token", "error", err)
	}
	return tokStr
}

func TestAuthenticatedMiddlewareHeaderSuccess(t *testing.T) {
	os.Setenv(utils.EnvJwtSignKey, "sekrit")
	os.Setenv(utils.EnvJwtExpireMinutes, "1")
	tokenString := createToken(t, "1234", 1)

	// create a handler to use as "next" which will verify the request
	nextHandler := createNextHandler(t, "1234")

	server := Server{}
	// create the handler to test, using our custom "next" handler
	handlerToTest := server.AuthenticatedMiddleware(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAuthenticatedMiddlewareHeaderExpiredToken(t *testing.T) {
	os.Setenv(utils.EnvJwtSignKey, "sekrit")
	os.Setenv(utils.EnvJwtExpireMinutes, "1")
	tokenString := createToken(t, "1234", -1)

	// create a handler to use as "next" which will verify the request
	nextHandler := createNextHandler(t, "1234")

	server := Server{}
	// create the handler to test, using our custom "next" handler
	handlerToTest := server.AuthenticatedMiddleware(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAuthenticatedMiddlewareHeaderWrongSigningKey(t *testing.T) {
	os.Setenv(utils.EnvJwtSignKey, "sekrit1")
	os.Setenv(utils.EnvJwtExpireMinutes, "1")
	tokenString := createToken(t, "1234", 1)

	// create a handler to use as "next" which will verify the request
	nextHandler := createNextHandler(t, "1234")

	server := Server{}
	// create the handler to test, using our custom "next" handler
	handlerToTest := server.AuthenticatedMiddleware(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAuthenticatedMiddlewareCookieSuccess(t *testing.T) {
	os.Setenv(utils.EnvJwtSignKey, "sekrit")
	os.Setenv(utils.EnvJwtExpireMinutes, "1")
	tokenString := createToken(t, "1234", 1)

	// create a handler to use as "next" which will verify the request
	nextHandler := createNextHandler(t, "1234")
	server := Server{}
	// create the handler to test, using our custom "next" handler
	handlerToTest := server.AuthenticatedMiddleware(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)
	cookie := createCookie(utils.AccessTokenCookieName, tokenString, time.Now().Add(time.Minute*time.Duration(1)))
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()
	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAuthenticatedMiddlewareCookieExpiredToken(t *testing.T) {
	os.Setenv(utils.EnvJwtSignKey, "sekrit")
	os.Setenv(utils.EnvJwtExpireMinutes, "1")
	tokenString := createToken(t, "1234", -1)

	// create a handler to use as "next" which will verify the request
	nextHandler := createNextHandler(t, "1234")
	server := Server{}
	// create the handler to test, using our custom "next" handler
	handlerToTest := server.AuthenticatedMiddleware(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)
	cookie := createCookie(utils.AccessTokenCookieName, tokenString, time.Now().Add(time.Minute*time.Duration(1)))
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()
	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAuthenticatedMiddlewareCookieWrongSigningKey(t *testing.T) {
	os.Setenv(utils.EnvJwtSignKey, "sekrit1")
	os.Setenv(utils.EnvJwtExpireMinutes, "1")
	tokenString := createToken(t, "1234", 1)

	// create a handler to use as "next" which will verify the request
	nextHandler := createNextHandler(t, "1234")
	server := Server{}
	// create the handler to test, using our custom "next" handler
	handlerToTest := server.AuthenticatedMiddleware(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)
	cookie := createCookie(utils.AccessTokenCookieName, tokenString, time.Now().Add(time.Minute*time.Duration(1)))
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()
	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func createNextHandler(t *testing.T, expectedSub string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value("sub")
		if val == nil {
			t.Error("sub not present")
		}
		valStr, ok := val.(string)
		if !ok {
			t.Error("not string")
		}
		if valStr != expectedSub {
			t.Fatalf("Got sub '%s', expected '%s'", valStr, expectedSub)
		}
	})

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
