package server

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"weight-tracker/internal/repository"
	"weight-tracker/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
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

func TestAuthenticatedMiddlewareCookieSuccess(t *testing.T) {
	os.Setenv(utils.EnvJwtSignKey, "sekrit")
	os.Setenv(utils.EnvJwtExpireMinutes, "1")
	tokenString := createToken(t, "1234", 1)

	nextHandler := createNextHandler(t, "1234")

	querierMock := querierMock{}
	querierMock.On("CheckIfTokenExists", mock.Anything, mock.Anything).Return(int64(0), sql.ErrNoRows)

	mockDb := &dbStub{
		repo: &querierMock,
	}

	server := Server{
		db: mockDb,
	}

	handlerToTest := server.AuthenticatedMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	cookie := createCookie(utils.AccessTokenCookieName, tokenString, time.Now().Add(time.Minute*time.Duration(1)))
	cookieRefresh := createCookie(utils.RefreshTokenCookieName, "refresh-token", time.Now().Add(time.Minute*time.Duration(30)))
	req.AddCookie(&cookie)
	req.AddCookie(&cookieRefresh)

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

	nextHandler := createNextHandler(t, "1234")

	querierMock := querierMock{}
	querierMock.On("CheckIfTokenExists", mock.Anything, mock.Anything).Return(int64(1), nil)

	mockDb := &dbStub{
		repo: &querierMock,
	}

	server := Server{
		db: mockDb,
	}

	handlerToTest := server.AuthenticatedMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	cookie := createCookie(utils.AccessTokenCookieName, tokenString, time.Now().Add(time.Minute*time.Duration(1)))
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()
	handlerToTest.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAuthenticatedMiddlewareCookieWrongSigningKey(t *testing.T) {
	os.Setenv(utils.EnvJwtSignKey, "sekrit1")
	os.Setenv(utils.EnvJwtExpireMinutes, "1")
	tokenString := createToken(t, "1234", 1)

	nextHandler := createNextHandler(t, "1234")

	querierMock := querierMock{}
	querierMock.On("CheckIfTokenExists", mock.Anything, mock.Anything).Return(int64(1), nil)

	mockDb := &dbStub{
		repo: &querierMock,
	}

	server := Server{
		db: mockDb,
	}

	handlerToTest := server.AuthenticatedMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	cookie := createCookie(utils.AccessTokenCookieName, tokenString, time.Now().Add(time.Minute*time.Duration(1)))
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()
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

		accessToken := r.Context().Value("access_token")
		if accessToken == nil {
			t.Error("access_token not present")
		}
		accessTokenStr, ok := accessToken.(string)
		if !ok {
			t.Error("access_token not string")
		}

		if accessTokenStr == "" {
			t.Error("access_token is empty")
		}

		refreshToken := r.Context().Value("refresh_token")
		if refreshToken == nil {
			t.Error("refresh_token not present")
		}
		refreshTokenStr, ok := refreshToken.(string)
		if !ok {
			t.Error("refresh_token not string")
		}
		if refreshTokenStr != "refresh-token" {
			t.Fatalf("Got refresh_token '%s', expected 'refresh-token'", refreshTokenStr)
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

type dbStub struct {
	repo repository.Querier
}

func (m *dbStub) Health() map[string]string {
	return map[string]string{}
}

func (m *dbStub) Close() error {
	return nil
}

func (m *dbStub) GetRepository() repository.Querier {
	return m.repo
}

type querierMock struct {
	mock.Mock
}

func (m *querierMock) CheckIfTokenExists(ctx context.Context, arg repository.CheckIfTokenExistsParams) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *querierMock) ReopenWorkoutById(ctx context.Context, arg repository.ReopenWorkoutByIdParams) (int64, error) {
	panic("not implemented")
}

func (m *querierMock) CompleteWorkoutById(ctx context.Context, arg repository.CompleteWorkoutByIdParams) (int64, error) {
	panic("not implemented")
}

func (m *querierMock) CreateExerciseAndReturnId(ctx context.Context, arg repository.CreateExerciseAndReturnIdParams) (string, error) {
	panic("not implemented")
}
func (m *querierMock) CreateExerciseItemAndReturnId(ctx context.Context, arg repository.CreateExerciseItemAndReturnIdParams) (string, error) {
	panic("not implemented")
}
func (m *querierMock) CreateExerciseTypeAndReturnId(ctx context.Context, arg repository.CreateExerciseTypeAndReturnIdParams) (string, error) {
	panic("not implemented")
}
func (m *querierMock) CreateExpiredToken(ctx context.Context, arg repository.CreateExpiredTokenParams) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) CreateSetAndReturnId(ctx context.Context, arg repository.CreateSetAndReturnIdParams) (string, error) {
	panic("not implemented")
}
func (m *querierMock) CreateUserAndReturnId(ctx context.Context, arg repository.CreateUserAndReturnIdParams) (string, error) {
	panic("not implemented")
}
func (m *querierMock) CreateWorkoutAndReturnId(ctx context.Context, arg repository.CreateWorkoutAndReturnIdParams) (string, error) {
	panic("not implemented")
}
func (m *querierMock) DeleteExerciseById(ctx context.Context, arg repository.DeleteExerciseByIdParams) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) DeleteExerciseItemById(ctx context.Context, arg repository.DeleteExerciseItemByIdParams) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) DeleteExerciseTypeById(ctx context.Context, arg repository.DeleteExerciseTypeByIdParams) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) DeleteExpiredTokens(ctx context.Context, currTime string) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) DeleteSetById(ctx context.Context, arg repository.DeleteSetByIdParams) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) DeleteUser(ctx context.Context, id string) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) DeleteWorkoutById(ctx context.Context, arg repository.DeleteWorkoutByIdParams) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) EmailExists(ctx context.Context, email any) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) GetAllExerciseTypes(ctx context.Context, userID string) ([]repository.ExerciseType, error) {
	panic("not implemented")
}
func (m *querierMock) GetAllExercises(ctx context.Context, userID string) ([]repository.Exercise, error) {
	panic("not implemented")
}
func (m *querierMock) GetAllSets(ctx context.Context, userID string) ([]repository.Set, error) {
	panic("not implemented")
}
func (m *querierMock) GetAllWorkouts(ctx context.Context, arg repository.GetAllWorkoutsParams) ([]repository.Workout, error) {
	panic("not implemented")
}
func (m *querierMock) GetAllWorkoutsCount(ctx context.Context, userID string) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) GetByEmail(ctx context.Context, email any) (repository.User, error) {
	panic("not implemented")
}
func (m *querierMock) GetByUserId(ctx context.Context, id string) (repository.User, error) {
	panic("not implemented")
}
func (m *querierMock) GetByUsername(ctx context.Context, username string) (repository.User, error) {
	panic("not implemented")
}
func (m *querierMock) GetExerciseById(ctx context.Context, arg repository.GetExerciseByIdParams) (repository.Exercise, error) {
	panic("not implemented")
}
func (m *querierMock) GetExerciseItemById(ctx context.Context, arg repository.GetExerciseItemByIdParams) (repository.ExerciseItem, error) {
	panic("not implemented")
}
func (m *querierMock) GetExerciseItemsByWorkoutId(ctx context.Context, arg repository.GetExerciseItemsByWorkoutIdParams) ([]repository.ExerciseItem, error) {
	panic("not implemented")
}
func (m *querierMock) GetExerciseTypeById(ctx context.Context, arg repository.GetExerciseTypeByIdParams) (repository.ExerciseType, error) {
	panic("not implemented")
}
func (m *querierMock) GetExercisesByWorkoutId(ctx context.Context, arg repository.GetExercisesByWorkoutIdParams) ([]repository.Exercise, error) {
	panic("not implemented")
}
func (m *querierMock) GetExercisesByExerciseItemId(ctx context.Context, arg repository.GetExercisesByExerciseItemIdParams) ([]repository.Exercise, error) {
	panic("not implemented")
}
func (m *querierMock) GetLastWeightRepsByExerciseTypeId(ctx context.Context, arg repository.GetLastWeightRepsByExerciseTypeIdParams) (repository.GetLastWeightRepsByExerciseTypeIdRow, error) {
	panic("not implemented")
}
func (m *querierMock) GetMaxWeightRepsByExerciseTypeId(ctx context.Context, arg repository.GetMaxWeightRepsByExerciseTypeIdParams) (repository.GetMaxWeightRepsByExerciseTypeIdRow, error) {
	panic("not implemented")
}
func (m *querierMock) GetSetById(ctx context.Context, arg repository.GetSetByIdParams) (repository.Set, error) {
	panic("not implemented")
}
func (m *querierMock) GetSetsByExerciseId(ctx context.Context, arg repository.GetSetsByExerciseIdParams) ([]repository.Set, error) {
	panic("not implemented")
}
func (m *querierMock) GetStatisticsBetweenDates(ctx context.Context, arg repository.GetStatisticsBetweenDatesParams) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) GetStatisticsSinceDate(ctx context.Context, arg repository.GetStatisticsSinceDateParams) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) GetUnverifiedUsers(ctx context.Context) ([]repository.User, error) {
	panic("not implemented")
}
func (m *querierMock) GetWorkoutById(ctx context.Context, arg repository.GetWorkoutByIdParams) (repository.Workout, error) {
	panic("not implemented")
}
func (m *querierMock) UpdateExerciseType(ctx context.Context, arg repository.UpdateExerciseTypeParams) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) UpdateExerciseItemType(ctx context.Context, arg repository.UpdateExerciseItemTypeParams) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) UpdateUser(ctx context.Context, arg repository.UpdateUserParams) (int64, error) {
	panic("not implemented")
}
func (m *querierMock) UpdateWorkoutById(ctx context.Context, arg repository.UpdateWorkoutByIdParams) (int64, error) {
	panic("not implemented")
}
