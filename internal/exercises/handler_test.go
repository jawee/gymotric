package exercises

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) GetAll(context context.Context, userId string) ([]Exercise, error) {
	args := s.Called(context, userId)
	return args.Get(0).([]Exercise), args.Error(1)
}
func (s *serviceMock) GetByWorkoutId(context context.Context, workoutId string, userId string) ([]Exercise, error) {
	args := s.Called(context, workoutId, userId)
	return args.Get(0).([]Exercise), args.Error(1)
}
func (s *serviceMock) DeleteById(context context.Context, id string, userId string) error {
	args := s.Called(context, id, userId)
	return args.Error(0)
}
func (s *serviceMock) CreateAndReturnId(context context.Context, exercise createExerciseRequest, workoutId string, userId string) (string, error) {
	args := s.Called(context, exercise, workoutId, userId)
	return args.String(0), args.Error(1)
}

func TestGetExercisesByWorkoutIdHandlerSuccess(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"
	exerciseTypeId := "exerciseTypeId"
	req, err := http.NewRequest("GET", "/workouts/"+workoutId+"/exercises", nil)
	req.SetPathValue("id", workoutId)

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sub", userId)
	req = req.WithContext(ctx)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetByWorkoutId", ctx, workoutId, userId).
		Return([]Exercise{
			{ID: "a", Name: "a", WorkoutID: workoutId, ExerciseTypeID: exerciseTypeId},
		}, nil).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getExercisesByWorkoutIdHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("content type header does not match: got %v want %v",
			ctype, "application/json")
	}

	expected := `{"exercises":[{"id":"a","name":"a","workout_id":"` + workoutId + `","exercise_type_id":"` + exerciseTypeId + `"}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetExercisesByWorkoutIdHandlerNotFound(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("GET", "/workouts/"+workoutId+"/exercises", nil)
	req.SetPathValue("id", workoutId)

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sub", userId)
	req = req.WithContext(ctx)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetByWorkoutId", ctx, workoutId, userId).
		Return([]Exercise{}, sql.ErrNoRows).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getExercisesByWorkoutIdHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	serviceMock.AssertExpectations(t)
}

func TestCreateExerciseHandler(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"
	exerciseTypeId := "exerciseTypeId"

	reqBodyObj := createExerciseRequest{
		ExerciseTypeID: exerciseTypeId,
	}

	reqBody, err := json.Marshal(reqBodyObj)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/workouts/"+workoutId+"/exercises", bytes.NewBuffer(reqBody))
	req.SetPathValue("id", workoutId)

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sub", userId)
	req = req.WithContext(ctx)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func (input createExerciseRequest) bool {
			return true
		}),  workoutId, userId).
		Return("abc", nil).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.createExerciseHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("content type header does not match: got %v want %v",
			ctype, "application/json")
	}

	expected := `{"id":"abc"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}
