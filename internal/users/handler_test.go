package users

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"weight-tracker/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (m *serviceMock) Register(ctx context.Context, arg registrationRequest) (string, error) {
	args := m.Called(ctx, arg)
	return args.String(0), args.Error(1)
}

func (m *serviceMock) ConfirmAccount(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *serviceMock) CreateAccountConfirmationToken(ctx context.Context, userId string) (string, error) {
	args := m.Called(ctx, userId)
	return args.String(0), args.Error(1)
}

func (m *serviceMock) CreateAndReturnId(ctx context.Context, arg createUserAndReturnIdRequest) (string, error) {
	args := m.Called(ctx, arg)
	return args.String(0), args.Error(1)
}

func (m *serviceMock) Login(ctx context.Context, arg loginRequest) (loginResponse, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(loginResponse), args.Error(1)
}

func (m *serviceMock) Logout(ctx context.Context, accessToken string, refreshToken string) error {
	args := m.Called(ctx, accessToken, refreshToken)
	return args.Error(0)
}

func (m *serviceMock) CreateToken(userId string) (string, error) {
	args := m.Called(userId)
	return args.String(0), args.Error(1)
}

func (m *serviceMock) GetByUserId(ctx context.Context, userId string) (getMeResponse, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(getMeResponse), args.Error(1)
}

func (m *serviceMock) ChangePassword(ctx context.Context, request changePasswordRequest, userId string) error {
	args := m.Called(ctx, request, userId)
	return args.Error(0)
}

func (m *serviceMock) CreateConfirmationToken(ctx context.Context, userId string, email string) (string, error) {
	args := m.Called(ctx, userId, email)
	return args.String(0), args.Error(1)
}

func (m *serviceMock) CreateResetPasswordToken(ctx context.Context, userId string) (string, error) {
	args := m.Called(ctx, userId)
	return args.String(0), args.Error(1)
}

func (m *serviceMock) ConfirmEmail(ctx context.Context, userId string, email string) error {
	args := m.Called(ctx, userId, email)
	return args.Error(0)
}

func (m *serviceMock) GetByEmail(ctx context.Context, email string) (getMeResponse, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(getMeResponse), args.Error(1)
}

func (m *serviceMock) ResetPassword(ctx context.Context, userId string, newPassword string) error {
	args := m.Called(ctx, userId, newPassword)
	return args.Error(0)
}

func (m *serviceMock) IsTokenValid(ctx context.Context, token string, tokenType string) bool {
	args := m.Called(ctx, token, tokenType)
	return args.Bool(0)
}

func populateContextWithSub(req *http.Request, userId string) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "sub", userId)
	return req.WithContext(ctx)
}

func populateContextWithSubAndTokens(req *http.Request, userId string, accessToken string, refreshToken string) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "sub", userId)
	ctx = context.WithValue(ctx, "access_token", accessToken)
	ctx = context.WithValue(ctx, "refresh_token", refreshToken)
	return req.WithContext(ctx)
}

var invalidJsonBytes = []byte("{invalid}")

func TestCreateUserHandler(t *testing.T) {
	jsonReqObj := createUserAndReturnIdRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	jsonReq, err := json.Marshal(jsonReqObj)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonReq))

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CreateAndReturnId", req.Context(), mock.MatchedBy(func(input createUserAndReturnIdRequest) bool {
		return true
	})).Return("userId", nil).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.createUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":"userId"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestCreateUserHandlerInvalidRequest(t *testing.T) {
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(invalidJsonBytes))

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.createUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

var testError = errors.New("Testerror")

func TestCreateUserHandlerServiceErr(t *testing.T) {
	jsonReqObj := createUserAndReturnIdRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	jsonReq, err := json.Marshal(jsonReqObj)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonReq))

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CreateAndReturnId", req.Context(), mock.MatchedBy(func(input createUserAndReturnIdRequest) bool {
		return true
	})).Return("", testError).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.createUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	serviceMock.AssertExpectations(t)
}

func TestLoginHandler(t *testing.T) {
	os.Setenv(utils.EnvJwtExpireMinutes, "10")
	os.Setenv(utils.EnvJwtSignKey, "test")
	os.Setenv(utils.EnvJwtRefreshExpireMinutes, "10")

	jsonReqObj := loginRequest{
		Username: "testuser",
		Password: "testpassword",
	}

	jsonReq, err := json.Marshal(jsonReqObj)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}
	serviceMock := serviceMock{}
	serviceMock.On("Login", req.Context(), mock.MatchedBy(func(input loginRequest) bool {
		return input.Username == "testuser" && input.Password == "testpassword"
	})).Return(loginResponse{
		Token:  "asdf",
		UserId: "userId",
	}, nil).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.loginHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	cookies := rr.Result().Cookies()
	assert.Condition(t, func() bool {
		for _, cookie := range cookies {
			if cookie.Name == utils.RefreshTokenCookieName {
				return true
			}
		}
		return false
	}, "handler did not set refresh cookie")

	assert.Condition(t, func() bool {
		for _, cookie := range cookies {
			if cookie.Name == utils.AccessTokenCookieName {
				return true
			}
		}
		return false
	}, "handler did not set access cookie")

	serviceMock.AssertExpectations(t)
}

func TestLoginHandlerJwtExpireNotParseable(t *testing.T) {
	os.Setenv(utils.EnvJwtExpireMinutes, "notparseable")

	jsonReqObj := loginRequest{
		Username: "testuser",
		Password: "testpassword",
	}

	jsonReq, err := json.Marshal(jsonReqObj)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}
	serviceMock := serviceMock{}

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.loginHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestLoginHandlerInvalidJson(t *testing.T) {
	os.Setenv(utils.EnvJwtExpireMinutes, "10")
	req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(invalidJsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	serviceMock := serviceMock{}

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.loginHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestLoginHandlerServiceErr(t *testing.T) {
	os.Setenv(utils.EnvJwtExpireMinutes, "10")

	jsonReqObj := loginRequest{
		Username: "testuser",
		Password: "testpassword",
	}

	jsonReq, err := json.Marshal(jsonReqObj)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}
	serviceMock := serviceMock{}
	serviceMock.On("Login", req.Context(), mock.MatchedBy(func(input loginRequest) bool {
		return input.Username == "testuser" && input.Password == "testpassword"
	})).Return(loginResponse{}, testError).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.loginHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetSubjectFromCookie(t *testing.T) {
	userId := "testuserId"
	os.Setenv(utils.EnvJwtSignKey, "testsigningkey")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * time.Duration(10))),
		Issuer:    "weight-tracker",
		Subject:   userId,
		Audience:  []string{"weight-tracker"},
	})
	signedToken, err := token.SignedString([]byte("testsigningkey"))

	if err != nil {
		t.Fatal(err)
	}

	cookies := []*http.Cookie{
		{
			Name:  utils.AccessTokenCookieName,
			Value: signedToken,
		},
	}
	serviceMock := serviceMock{}
	serviceMock.On("IsTokenValid", mock.Anything, signedToken, "access_token").Return(true).Once()
	s := handler{service: &serviceMock}

	sub, err := s.getSubjectFromCookie(context.Background(), utils.AccessTokenCookieName, "testsigningkey", "access_token", cookies) 

	assert.Nil(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, userId, sub)
}

func TestGetSubjectFromCookieNoCookieFound(t *testing.T) {
	cookies := []*http.Cookie{}
	serviceMock := serviceMock{}
	s := handler{service: &serviceMock}
	sub, err := s.getSubjectFromCookie(context.Background(), utils.AccessTokenCookieName, "testsigningkey", "access_token", cookies)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, NoTokenFoundError)
	assert.Empty(t, sub)
}

func TestMeHandler(t *testing.T) {
	userId := "testuserId"
	req, err := http.NewRequest("GET", "/users/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = populateContextWithSub(req, userId)

	serviceMock := serviceMock{}
	serviceMock.On("GetByUserId", req.Context(), userId).Return(getMeResponse{
		ID:        userId,
		Username:  "testuser",
		Email:     "testuser@email.com",
		CreatedOn: "2025-04-19T08:16:15Z",
		UpdatedOn: "2025-04-19T08:16:15Z",
	}, nil).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.meHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":"testuserId","username":"testuser","email":"testuser@email.com","created_on":"2025-04-19T08:16:15Z","updated_on":"2025-04-19T08:16:15Z"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestMeHandlerServiceErr(t *testing.T) {
	userId := "testuserId"
	req, err := http.NewRequest("GET", "/users/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = populateContextWithSub(req, userId)

	serviceMock := serviceMock{}
	serviceMock.On("GetByUserId", req.Context(), userId).Return(getMeResponse{}, testError).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.meHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	serviceMock.AssertExpectations(t)
}

func TestLogoutHandler(t *testing.T) {
	userId := "testuserId"
	accessToken := "testAccessToken"
	refreshToken := "testRefreshToken"

	req, err := http.NewRequest("GET", "/users/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = populateContextWithSubAndTokens(req, userId, accessToken, refreshToken)

	serviceMock := serviceMock{}
	serviceMock.On("Logout", req.Context(), accessToken, refreshToken).Return(nil).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.logoutHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
	cookies := rr.Result().Cookies()
	assert.Condition(t, func() bool {
		for _, cookie := range cookies {
			if cookie.Name == utils.AccessTokenCookieName && cookie.Value == "" && cookie.Expires.Before(time.Now().Add(time.Second)) {
				return true
			}
		}
		return false
	}, "handler did not set access cookie")
	assert.Condition(t, func() bool {
		for _, cookie := range cookies {
			if cookie.Name == utils.RefreshTokenCookieName && cookie.Value == "" && cookie.Expires.Before(time.Now().Add(time.Second)) {
				return true
			}
		}
		return false
	}, "handler did not set refresh cookie")
}
