package workouts

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

func populateContextWithSub(req *http.Request, userId string) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "sub", userId)
	return req.WithContext(ctx)
}

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

func TestGetAllWorkoutsHandlerGetAllErr(t *testing.T) {
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

func TestGetWorkoutByIdHandler(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("GET", "/workouts/"+workoutId, nil)
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetById", req.Context(), workoutId, userId).
		Return(Workout{
			ID:   workoutId,
			Name: "workoutName",
		}, nil).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getWorkoutByIdHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetWorkoutByIdHandlerNotFound(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("GET", "/workouts/"+workoutId, nil)
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetById", req.Context(), workoutId, userId).
		Return(Workout{}, sql.ErrNoRows).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getWorkoutByIdHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	serviceMock.AssertExpectations(t)
}

func TestGetWorkoutByIdHandlerServiceErr(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("GET", "/workouts/"+workoutId, nil)
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("GetById", req.Context(), workoutId, userId).
		Return(Workout{}, testError).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.getWorkoutByIdHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	serviceMock.AssertExpectations(t)
}

func TestCreateWorkoutHandler(t *testing.T) {
	userId := "userId"

	reqBodyObj := createWorkoutRequest{
		Name: "workoutName",
	}

	reqBody, err := json.Marshal(reqBodyObj)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/workouts", bytes.NewBuffer(reqBody))

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CreateAndReturnId", req.Context(), mock.MatchedBy(func(input createWorkoutRequest) bool {
		return input.Name == reqBodyObj.Name
	}), userId).Return("id", nil).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.createWorkoutHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := `{"id":"id"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestCreateWorkoutHandlerServiceErr(t *testing.T) {
	userId := "userId"

	reqBodyObj := createWorkoutRequest{
		Name: "workoutName",
	}

	reqBody, err := json.Marshal(reqBodyObj)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/workouts", bytes.NewBuffer(reqBody))

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CreateAndReturnId", req.Context(), mock.Anything, userId).Return("", testError).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.createWorkoutHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	serviceMock.AssertExpectations(t)
}

func TestCompleteWorkoutByIdHandler(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("POST", "/workouts/"+workoutId+"/complete", nil)
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CompleteById", req.Context(), workoutId, userId).Return(nil).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.completeWorkoutById)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	serviceMock.AssertExpectations(t)
}

func TestCompleteWorkoutByIdHandlerServiceErr(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("POST", "/workouts/"+workoutId+"/complete", nil)
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CompleteById", req.Context(), workoutId, userId).Return(testError).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.completeWorkoutById)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	serviceMock.AssertExpectations(t)
}

func TestUpdateWorkoutByIdHandler(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	reqBodyObj := updateWorkoutRequest{
		Note: "workoutNote",
	}

	reqBody, err := json.Marshal(reqBodyObj)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/workouts/"+workoutId, bytes.NewBuffer(reqBody))
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("UpdateWorkoutById", req.Context(), workoutId, mock.MatchedBy(func(input updateWorkoutRequest) bool {
		return input.Note == reqBodyObj.Note
	}), userId).Return(nil).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.updateWorkoutByIdHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	serviceMock.AssertExpectations(t)
}

var invalidJsonBytes = []byte("{invalidJson}")

func TestUpdateWorkoutByIdHandlerJsonErr(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("PUT", "/workouts/"+workoutId, bytes.NewBuffer(invalidJsonBytes))
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.updateWorkoutByIdHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestUpdateWorkoutByIdHandlerServiceErr(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	reqBodyObj := updateWorkoutRequest{
		Note: "workoutNote",
	}

	reqBody, err := json.Marshal(reqBodyObj)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/workouts/"+workoutId, bytes.NewBuffer(reqBody))
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("UpdateWorkoutById", req.Context(), workoutId, mock.Anything, userId).Return(testError).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.updateWorkoutByIdHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	serviceMock.AssertExpectations(t)
}

func TestDeleteWorkoutByIdHandler(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("DELETE", "/workouts/"+workoutId, nil)
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("DeleteById", req.Context(), workoutId, userId).Return(nil).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.deleteWorkoutByIdHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	serviceMock.AssertExpectations(t)
}

func TestDeleteWorkoutByIdHandlerServiceErr(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("DELETE", "/workouts/"+workoutId, nil)
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("DeleteById", req.Context(), workoutId, userId).Return(testError).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.deleteWorkoutByIdHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	serviceMock.AssertExpectations(t)
}

func TestCloneWorkoutByIdHandler(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("POST", "/workouts/"+workoutId+"/clone", nil)
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CloneByIdAndReturnId", req.Context(), workoutId, userId).Return("newWorkoutId", nil).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.cloneWorkoutById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := `{"id":"newWorkoutId"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	serviceMock.AssertExpectations(t)
}

func TestCloneWorkoutByIdHandlerServiceErr(t *testing.T) {
	userId := "userId"
	workoutId := "workoutId"

	req, err := http.NewRequest("POST", "/workouts/"+workoutId+"/clone", nil)
	req.SetPathValue("id", workoutId)

	req = populateContextWithSub(req, userId)

	if err != nil {
		t.Fatal(err)
	}

	serviceMock := serviceMock{}
	serviceMock.On("CloneByIdAndReturnId", req.Context(), workoutId, userId).Return("", testError).Once()

	rr := httptest.NewRecorder()
	s := handler{service: &serviceMock}
	handler := http.HandlerFunc(s.cloneWorkoutById)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	serviceMock.AssertExpectations(t)
}
