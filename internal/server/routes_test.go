package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"weight-tracker/internal/repository"
)

// // database.Service
// CreateExerciseTypeAndReturnId
type MockDatabaseService struct{}

func (m MockDatabaseService) Health() map[string]string {
	return map[string]string{}
}

func (m MockDatabaseService) Close() error {
	return nil
}

func (m MockDatabaseService) GetRepository() repository.Querier {
	return MockRepository{}
}

type MockRepository struct{}

// CompleteWorkoutById implements repository.Querier.
func (m MockRepository) CompleteWorkoutById(ctx context.Context, arg repository.CompleteWorkoutByIdParams) error {
	panic("unimplemented")
}

// CreateExerciseAndReturnId implements repository.Querier.
func (m MockRepository) CreateExerciseAndReturnId(ctx context.Context, arg repository.CreateExerciseAndReturnIdParams) (string, error) {
	panic("unimplemented")
}

// CreateExerciseTypeAndReturnId implements repository.Querier.
func (m MockRepository) CreateExerciseTypeAndReturnId(ctx context.Context, arg repository.CreateExerciseTypeAndReturnIdParams) (string, error) {
	panic("unimplemented")
}

// CreateSetAndReturnId implements repository.Querier.
func (m MockRepository) CreateSetAndReturnId(ctx context.Context, arg repository.CreateSetAndReturnIdParams) (string, error) {
	panic("unimplemented")
}

// CreateWorkoutAndReturnId implements repository.Querier.
func (m MockRepository) CreateWorkoutAndReturnId(ctx context.Context, arg repository.CreateWorkoutAndReturnIdParams) (string, error) {
	panic("unimplemented")
}

// DeleteExerciseById implements repository.Querier.
func (m MockRepository) DeleteExerciseById(ctx context.Context, id string) error {
	panic("unimplemented")
}

// DeleteExerciseTypeById implements repository.Querier.
func (m MockRepository) DeleteExerciseTypeById(ctx context.Context, id string) error {
	panic("unimplemented")
}

// DeleteSetById implements repository.Querier.
func (m MockRepository) DeleteSetById(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetAllExerciseTypes implements repository.Querier.
func (m MockRepository) GetAllExerciseTypes(ctx context.Context) ([]repository.ExerciseType, error) {
	panic("unimplemented")
}

// GetAllExercises implements repository.Querier.
func (m MockRepository) GetAllExercises(ctx context.Context) ([]repository.Exercise, error) {
	panic("unimplemented")
}

// GetAllSets implements repository.Querier.
func (m MockRepository) GetAllSets(ctx context.Context) ([]repository.Set, error) {
	panic("unimplemented")
}

// GetAllWorkouts implements repository.Querier.
func (m MockRepository) GetAllWorkouts(ctx context.Context) ([]repository.Workout, error) {
	workouts := []repository.Workout{}
	return workouts, nil
}

// GetExerciseById implements repository.Querier.
func (m MockRepository) GetExerciseById(ctx context.Context, id string) (repository.Exercise, error) {
	panic("unimplemented")
}

// GetExerciseTypeById implements repository.Querier.
func (m MockRepository) GetExerciseTypeById(ctx context.Context, id string) (repository.ExerciseType, error) {
	panic("unimplemented")
}

// GetExercisesByWorkoutId implements repository.Querier.
func (m MockRepository) GetExercisesByWorkoutId(ctx context.Context, workoutID string) ([]repository.Exercise, error) {
	panic("unimplemented")
}

// GetSetById implements repository.Querier.
func (m MockRepository) GetSetById(ctx context.Context, id string) (repository.Set, error) {
	panic("unimplemented")
}

// GetSetsByExerciseId implements repository.Querier.
func (m MockRepository) GetSetsByExerciseId(ctx context.Context, exerciseID string) ([]repository.Set, error) {
	panic("unimplemented")
}

// GetWorkoutById implements repository.Querier.
func (m MockRepository) GetWorkoutById(ctx context.Context, id string) (repository.Workout, error) {
	panic("unimplemented")
}

func TestHandler(t *testing.T) {
	s := &Server{
		db: MockDatabaseService{},
	}
	server := httptest.NewServer(http.HandlerFunc(s.getAllWorkoutsHandler))
	defer server.Close()
	resp, err := http.Get(server.URL + "/workouts")
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"workouts\":[]}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}
