package workouts

import (
	"context"
	"fmt"
	"log/slog"
	"weight-tracker/internal/repository"
)

type Workout struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CompletedOn any    `json:"completed_on"`
	CreatedOn   string `json:"created_on"`
	UpdatedOn   string `json:"updated_on"`
	Note        string `json:"note"`
}

type WorkoutsRepository interface {
	CompleteById(ctx context.Context, arg repository.CompleteWorkoutByIdParams) (int64, error)
	CreateAndReturnId(ctx context.Context, arg repository.CreateWorkoutAndReturnIdParams) (string, error)
	GetAll(ctx context.Context, arg repository.GetAllWorkoutsParams) ([]Workout, error)
	GetAllCount(ctx context.Context, userId string) (int64, error)
	GetById(ctx context.Context, arg repository.GetWorkoutByIdParams) (Workout, error)
	DeleteById(ctx context.Context, arg repository.DeleteWorkoutByIdParams) error
	UpdateById(context context.Context, arg repository.UpdateWorkoutByIdParams) error
}

type workoutsRepository struct {
	repo repository.Querier
}

func (w *workoutsRepository) GetAllCount(ctx context.Context, userId string) (int64, error) {
	count, err := w.repo.GetAllWorkoutsCount(ctx, userId)
	if err != nil {
		return 0, fmt.Errorf("failed to get workouts count: %w", err)
	}

	return count, nil
}

func (w *workoutsRepository) UpdateById(ctx context.Context, arg repository.UpdateWorkoutByIdParams) error {
	rows, err := w.repo.UpdateWorkoutById(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to update workout: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("workout not found")
	}
	return nil
}

func (w *workoutsRepository) DeleteById(ctx context.Context, arg repository.DeleteWorkoutByIdParams) error {
	rows, err := w.repo.DeleteWorkoutById(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to delete workout: %w", err)
	}

	if rows == 0 {
		slog.Warn("Tried to delete workout that did not exist", "workoutId", arg.ID)
	}
	return nil
}

func (w *workoutsRepository) CompleteById(ctx context.Context, arg repository.CompleteWorkoutByIdParams) (int64, error) {
	return w.repo.CompleteWorkoutById(ctx, arg)
}

func (w *workoutsRepository) CreateAndReturnId(ctx context.Context, arg repository.CreateWorkoutAndReturnIdParams) (string, error) {
	return w.repo.CreateWorkoutAndReturnId(ctx, arg)
}

func (w *workoutsRepository) GetAll(ctx context.Context, arg repository.GetAllWorkoutsParams) ([]Workout, error) {
	workouts, err := w.repo.GetAllWorkouts(ctx, arg)
	if err != nil {
		return []Workout{}, fmt.Errorf("failed to get all workouts: %w", err)
	}

	result := []Workout{}
	for _, v := range workouts {
		result = append(result, newWorkout(v))
	}

	return result, nil
}

func newWorkout(v repository.Workout) Workout {
	workout := Workout{
		ID:          v.ID,
		Name:        v.Name,
		CompletedOn: v.CompletedOn,
		CreatedOn:   v.CreatedOn,
		UpdatedOn:   v.UpdatedOn,
	}

	if v.Note != nil {
		workout.Note = v.Note.(string)
	}

	return workout
}

func (w *workoutsRepository) GetById(ctx context.Context, arg repository.GetWorkoutByIdParams) (Workout, error) {
	workout, err := w.repo.GetWorkoutById(ctx, arg)
	if err != nil {
		return Workout{}, fmt.Errorf("failed to get workout by ID: %w", err)
	}

	result := newWorkout(workout)
	return result, nil
}
