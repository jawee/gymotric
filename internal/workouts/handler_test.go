package workouts

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) GetAll(context context.Context, userId string, page int, pageSize int) ([]Workout, error) {
	args := s.Called(context, userId, page, pageSize)
	return args.Get(0).([]Workout), args.Error(1)

}
func (s *serviceMock) GetAllCount(context context.Context, userId string) (int, error) {
	args := s.Called(context, userId)
	return args.Int(0), args.Error(1)
}
func (s *serviceMock) GetById(context context.Context, id string, userId string) (Workout, error) {
	args := s.Called(context, id, userId)
	return args.Get(0).(Workout), args.Error(1)
}
func (s *serviceMock) CreateAndReturnId(context context.Context, t createWorkoutRequest, userId string) (string, error) {
	args := s.Called(context, t, userId)
	return args.String(0), args.Error(1)
}
func (s *serviceMock) CompleteById(context context.Context, workoutId string, userId string) error {
	args := s.Called(context, workoutId, userId)
	return args.Error(0)
}
func (s *serviceMock) DeleteById(context context.Context, workoutId string, userId string) error {
	args := s.Called(context, workoutId, userId)
	return args.Error(0)
}
func (s *serviceMock) CloneByIdAndReturnId(context context.Context, workoutId string, userId string) (string, error) {
	args := s.Called(context, workoutId, userId)
	return args.String(0), args.Error(1)
}
func (s *serviceMock) UpdateWorkoutById(context context.Context, workoutId string, t updateWorkoutRequest, userId string) error {
	args := s.Called(context, workoutId, t, userId)
	return args.Error(0)
}

func TestGetAllWorkoutsHandler(t *testing.T) {
	userId := "userId"

	req, err := http.NewRequest("GET", "/workouts", nil)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetAll", req.Context(), userId, 1, 10).
		Return([]Workout{
			{
				ID:   "1",
				Name: "workoutName",
			},
		}, nil).Once()
	serviceMock.On("GetAllCount", req.Context(), userId).Return(1, nil).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getAllWorkoutsHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"data":[{"id":"1","name":"workoutName","completed_on":null,"created_on":"","updated_on":"","note":""}],"page":1,"page_size":10,"total":1,"total_pages":1}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetAllWorkoutsHandlerNoFound(t *testing.T) {
	userId := "userId"

	req, err := http.NewRequest("GET", "/workouts", nil)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetAll", req.Context(), userId, 1, 10).
		Return([]Workout{}, sql.ErrNoRows).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getAllWorkoutsHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetAllWorkoutsHandlerGetAllErr(t * testing.T) {
	userId := "userId"

	req, err := http.NewRequest("GET", "/workouts", nil)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetAll", req.Context(), userId, 1, 10).
		Return([]Workout{}, testError).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getAllWorkoutsHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetAllWorkoutsHandlerGetAllCountErr(t *testing.T) {
	userId := "userId"

	req, err := http.NewRequest("GET", "/workouts", nil)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetAll", req.Context(), userId, 1, 10).
		Return([]Workout{
			{
				ID:   "1",
				Name: "workoutName",
			},
		}, nil).Once()
	serviceMock.On("GetAllCount", req.Context(), userId).Return(0, testError).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getAllWorkoutsHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	serviceMock.AssertExpectations(t)
}

func populateContextWithSub(req *http.Request, userId string) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "sub", userId)
	return req.WithContext(ctx)
}
