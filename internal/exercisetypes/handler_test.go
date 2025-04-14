package exercisetypes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) GetAll(context context.Context, userId string) ([]ExerciseType, error) {
	args := s.Called(context, userId)
	return args.Get(0).([]ExerciseType), args.Error(1)
}
func (s *serviceMock) DeleteById(context context.Context, exerciseTypeId string, userId string) error {
	args := s.Called(context, exerciseTypeId, userId)
	return args.Error(0)
}
func (s *serviceMock) CreateAndReturnId(context context.Context, exerciseType createExerciseTypeRequest, userId string) (string, error) {
	args := s.Called(context, exerciseType, userId)
	return args.String(0), args.Error(1)
}
func (s *serviceMock) GetLastWeightRepsByExerciseTypeId(context context.Context, exerciseTypeId string, userId string) (MaxLastWeightReps, error) {
	args := s.Called(context, exerciseTypeId, userId)
	return args.Get(0).(MaxLastWeightReps), args.Error(1)
}
func (s *serviceMock) GetMaxWeightRepsByExerciseTypeId(context context.Context, exerciseTypeId string, userId string) (MaxLastWeightReps, error) {
	args := s.Called(context, exerciseTypeId, userId)
	return args.Get(0).(MaxLastWeightReps), args.Error(1)
}

func populateContextWithSub(req *http.Request, userId string) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "sub", userId)
	return req.WithContext(ctx)
}

func TestCreateExerciseTypeHandler(t *testing.T) {
	userId := "userId"

	reqBodyObj := createExerciseTypeRequest{
		Name: "exerciseName",
	}

	reqBody, err := json.Marshal(reqBodyObj)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/exercise-types", bytes.NewBuffer(reqBody))

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CreateAndReturnId", req.Context(), mock.MatchedBy(func(input createExerciseTypeRequest) bool {
		return input.Name == "exerciseName"
	}), userId).
		Return("abc", nil).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.createExerciseTypeHandler)

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
func TestCreateExerciseTypeHandlerErr(t *testing.T) {
	userId := "userId"

	reqBodyObj := createExerciseTypeRequest{
		Name: "exerciseName",
	}

	reqBody, err := json.Marshal(reqBodyObj)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/exercise-types", bytes.NewBuffer(reqBody))

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CreateAndReturnId", req.Context(), mock.MatchedBy(func(input createExerciseTypeRequest) bool {
		return true
	}), userId).
		Return("", errors.New("Failed")).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.createExerciseTypeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Failed to create exercise\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetAllExerciseTypesHandler(t *testing.T) {
	userId := "userId"

	req, err := http.NewRequest("GET", "/exercise-types", nil)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetAll", req.Context(), userId).
		Return([]ExerciseType{
			{
				ID:   "1",
				Name: "exerciseName",
			},
		}, nil).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getAllWorkoutTypesHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"exercise_types":[{"id":"1","name":"exerciseName"}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)

}
func TestGetAllExerciseTypesHandlerErr(t *testing.T) {
	userId := "userId"

	req, err := http.NewRequest("GET", "/exercise-types", nil)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetAll", req.Context(), userId).
		Return([]ExerciseType{}, errors.New("Err")).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getAllWorkoutTypesHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Body.String() != "\n" {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), "")
	}

	serviceMock.AssertExpectations(t)
}

func TestDeleteExerciseTypeByIdHandler(t *testing.T) {
	userId := "userId"
	exerciseTypeId := "exerciseTypeId"

	req, err := http.NewRequest("DELETE", "/exercise-types/"+exerciseTypeId, nil)
	req.SetPathValue("id", exerciseTypeId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("DeleteById", req.Context(), exerciseTypeId, userId).
		Return(nil).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.deleteExerciseTypeByIdHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	serviceMock.AssertExpectations(t)
}

func TestDeleteExerciseTypeByIdHandlerErr(t *testing.T) {
	userId := "userId"
	exerciseTypeId := "exerciseTypeId"

	req, err := http.NewRequest("DELETE", "/exercise-types/"+exerciseTypeId, nil)
	req.SetPathValue("id", exerciseTypeId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("DeleteById", req.Context(), exerciseTypeId, userId).
		Return(errors.New("Failed")).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.deleteExerciseTypeByIdHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Failed to delete exercise type\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}
