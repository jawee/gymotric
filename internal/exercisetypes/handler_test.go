package exercisetypes

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

var testError = errors.New("Testerror")

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
func (s *serviceMock) UpdateById(context context.Context, exerciseTypeId string, updateExerciseTypeRequest updateExerciseTypeRequest, userId string) error {
	args := s.Called(context, exerciseTypeId, updateExerciseTypeRequest, userId)
	return args.Error(0)
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
		Return("", testError).
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

	expected := `{"data":[{"id":"1","name":"exerciseName"}]}`
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
		Return([]ExerciseType{}, testError).
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
		Return(testError).
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

func TestGetLastSet(t *testing.T) {
	userId := "userId"
	exerciseTypeId := "exerciseTypeId"

	req, err := http.NewRequest("GET", "/exercise-types/"+exerciseTypeId+"/last", nil)
	req.SetPathValue("id", exerciseTypeId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetLastWeightRepsByExerciseTypeId", req.Context(), exerciseTypeId, userId).
		Return(MaxLastWeightReps{Weight: 100, Reps: 10}, nil).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getLastSet)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"data":{"weight":100,"reps":10}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetLastSetErr(t *testing.T) {
	userId := "userId"
	exerciseTypeId := "exerciseTypeId"

	req, err := http.NewRequest("GET", "/exercise-types/"+exerciseTypeId+"/last", nil)
	req.SetPathValue("id", exerciseTypeId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetLastWeightRepsByExerciseTypeId", req.Context(), exerciseTypeId, userId).
		Return(MaxLastWeightReps{}, testError).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getLastSet)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Body.String() != "\n" {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), "\n")
	}

	serviceMock.AssertExpectations(t)
}

func TestGetLastSetSqlErr(t *testing.T) {
	userId := "userId"
	exerciseTypeId := "exerciseTypeId"

	req, err := http.NewRequest("GET", "/exercise-types/"+exerciseTypeId+"/last", nil)
	req.SetPathValue("id", exerciseTypeId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetLastWeightRepsByExerciseTypeId", req.Context(), exerciseTypeId, userId).
		Return(MaxLastWeightReps{}, sql.ErrNoRows).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getLastSet)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Body.String() != `{"data":{"weight":0,"reps":0}}` {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), `{"data":{"weight":0,"reps":0}}`)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetMaxSet(t *testing.T) {
	userId := "userId"
	exerciseTypeId := "exerciseTypeId"

	req, err := http.NewRequest("GET", "/exercise-types/"+exerciseTypeId+"/max", nil)
	req.SetPathValue("id", exerciseTypeId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetMaxWeightRepsByExerciseTypeId", req.Context(), exerciseTypeId, userId).
		Return(MaxLastWeightReps{Weight: 110, Reps: 10}, nil).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getMaxSet)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"data":{"weight":110,"reps":10}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetMaxSetErr(t *testing.T) {
	userId := "userId"
	exerciseTypeId := "exerciseTypeId"

	req, err := http.NewRequest("GET", "/exercise-types/"+exerciseTypeId+"/max", nil)
	req.SetPathValue("id", exerciseTypeId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetMaxWeightRepsByExerciseTypeId", req.Context(), exerciseTypeId, userId).
		Return(MaxLastWeightReps{}, testError).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getMaxSet)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Body.String() != "\n" {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), "\n")
	}

	serviceMock.AssertExpectations(t)
}

func TestGetMaxSetSqlErr(t *testing.T) {
	userId := "userId"
	exerciseTypeId := "exerciseTypeId"

	req, err := http.NewRequest("GET", "/exercise-types/"+exerciseTypeId+"/max", nil)
	req.SetPathValue("id", exerciseTypeId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetMaxWeightRepsByExerciseTypeId", req.Context(), exerciseTypeId, userId).
		Return(MaxLastWeightReps{}, sql.ErrNoRows).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getMaxSet)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Body.String() != `{"data":{"weight":0,"reps":0}}` {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), `{"data":{"weight":0,"reps":0}}`)
	}

	serviceMock.AssertExpectations(t)
}

func TestUpdateExerciseTypeByIdHandler(t *testing.T) {
	userId := "userId"
	exerciseTypeId := "exerciseTypeId"

	reqBodyObj := updateExerciseTypeRequest{
		Name: "exerciseName",
	}

	reqBody, err := json.Marshal(reqBodyObj)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/exercise-types/"+exerciseTypeId, bytes.NewBuffer(reqBody))
	req.SetPathValue("id", exerciseTypeId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("UpdateById", req.Context(), exerciseTypeId, mock.MatchedBy(func(input updateExerciseTypeRequest) bool {
		return input.Name == "exerciseName"
	}), userId).
		Return(nil).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.updateExerciseTypeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	serviceMock.AssertExpectations(t)
}
