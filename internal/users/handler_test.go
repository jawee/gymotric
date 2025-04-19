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
	"weight-tracker/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (m *serviceMock) CreateAndReturnId(ctx context.Context, arg createUserAndReturnIdRequest) (string, error) {
	args := m.Called(ctx, arg)
	return args.String(0), args.Error(1)
}

func (m *serviceMock) Login(ctx context.Context, arg loginRequest) (loginResponse, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(loginResponse), args.Error(1)
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

func populateContextWithSub(req *http.Request, userId string) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "sub", userId)
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

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response["access_token"] != "asdf" {
		t.Errorf("handler returned unexpected access_token: got %v want %v", response["access_token"], "asdf")
	}
	if response["token_type"] != "Bearer" {
		t.Errorf("handler returned unexpected token_type: got %v want %v", response["token_type"], "Bearer")
	}

	if val, ok := response["expires_in"].(float64); !ok || val != 600 {
		t.Errorf("handler returned unexpected expires_in: got %v want %v", val, 600)
	}

	if _, ok := response["refresh_token"]; !ok {
		t.Errorf("handler returned unexpected refresh_token: got %v want %v", response["refresh_token"], "asdf")
	}

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
