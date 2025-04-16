package sets

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

func (s *serviceMock) GetByExerciseId(context context.Context, exerciseId string, userId string) ([]Set, error) {
	args := s.Called(context, exerciseId, userId)
	return args.Get(0).([]Set), args.Error(1)
}
func (s *serviceMock) DeleteById(context context.Context, setId string, userId string) error {
	args := s.Called(context, setId, userId)
	return args.Error(0)
}
func (s *serviceMock) CreateAndReturnId(context context.Context, t createSetRequest, exerciseId string, userId string) (string, error) {
	args := s.Called(context, t, exerciseId, userId)
	return args.String(0), args.Error(1)
}

func populateContextWithSub(req *http.Request, userId string) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "sub", userId)
	return req.WithContext(ctx)
}

func TestDeleteSetByIdHandler(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"
	exerciseId := "exerciseId"
	setId := "setId"

	req, err := http.NewRequest("DELETE", "/workouts/"+workoutId+"/exercises/"+exerciseId+"/sets/"+setId, nil)
	req.SetPathValue("id", workoutId)
	req.SetPathValue("exerciseId", exerciseId)
	req.SetPathValue("setId", setId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("DeleteById", req.Context(), setId, userId).
		Return(nil).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.deleteSetByIdHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	serviceMock.AssertExpectations(t)
}

func TestDeleteSetByIdHandlerErr(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"
	exerciseId := "exerciseId"
	setId := "setId"

	req, err := http.NewRequest("DELETE", "/workouts/"+workoutId+"/exercises/"+exerciseId+"/sets/"+setId, nil)
	req.SetPathValue("id", workoutId)
	req.SetPathValue("exerciseId", exerciseId)
	req.SetPathValue("setId", setId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("DeleteById", req.Context(), setId, userId).
		Return(errors.New("Failed")).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.deleteSetByIdHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Failed to delete set\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetSetsByExerciseIdHandler(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"
	exerciseId := "exerciseId"

	req, err := http.NewRequest("GET", "/workouts/"+workoutId+"/exercises/"+exerciseId+"/sets", nil)
	req.SetPathValue("id", workoutId)
	req.SetPathValue("exerciseId", exerciseId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetByExerciseId", req.Context(), exerciseId, userId).
		Return([]Set{{ID: "a", Repetitions: 10, Weight: 100, ExerciseID: exerciseId}}, nil).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getSetsByExerciseIdHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"sets":[{"id":"a","repetitions":10,"weight":100,"exercise_id":"exerciseId"}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetSetsByExerciseIdHandlerNoExerciseId(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("GET", "/workouts/"+workoutId+"/exercises/sets", nil)
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getSetsByExerciseIdHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := "Failed to get sets\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetSetsByExerciseIdHandlerServiceErr(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"
	exerciseId := "exerciseId"

	req, err := http.NewRequest("GET", "/workouts/"+workoutId+"/exercises/"+exerciseId+"/sets", nil)
	req.SetPathValue("id", workoutId)
	req.SetPathValue("exerciseId", exerciseId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetByExerciseId", req.Context(), exerciseId, userId).
		Return([]Set{}, errors.New("Failed")).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getSetsByExerciseIdHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Failed to get sets\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestCreateSetHandler(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"
	exerciseId := "exerciseId"

	reqBodyObj := createSetRequest{
		Repetitions: 10,
		Weight:      100,
	}

	reqBody, err := json.Marshal(reqBodyObj)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/workouts/"+workoutId+"/exercises/"+exerciseId+"/sets", bytes.NewBuffer(reqBody))
	req.SetPathValue("id", workoutId)
	req.SetPathValue("exerciseId", exerciseId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CreateAndReturnId", req.Context(), mock.MatchedBy(func(input createSetRequest) bool {
		return input.Repetitions == 10 && input.Weight == 100
	}), exerciseId, userId).
		Return("setId", nil).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.createSetHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":"setId"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestCreateSetHandlerErr(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"
	exerciseId := "exerciseId"

	reqBodyObj := createSetRequest{
		Repetitions: 10,
		Weight:      100,
	}

	reqBody, err := json.Marshal(reqBodyObj)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/workouts/"+workoutId+"/exercises/"+exerciseId+"/sets", bytes.NewBuffer(reqBody))
	req.SetPathValue("id", workoutId)
	req.SetPathValue("exerciseId", exerciseId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CreateAndReturnId", req.Context(), mock.MatchedBy(func(input createSetRequest) bool {
		return input.Repetitions == 10 && input.Weight == 100
	}), exerciseId, userId).
		Return("", errors.New("Failed")).
		Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.createSetHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Failed to create set\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}
