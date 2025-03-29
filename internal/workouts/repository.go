package workouts

import (
	"context"
	"fmt"
	"log/slog"
	"weight-tracker/internal/repository"
)

type Workout struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	CompletedOn interface{} `json:"completed_on"`
	CreatedOn   string      `json:"created_on"`
	UpdatedOn   string      `json:"updated_on"`
	Note        string      `json:"note"`
}

type WorkoutsRepository interface {
	CompleteById(ctx context.Context, arg repository.CompleteWorkoutByIdParams) (int64, error)
	CreateAndReturnId(ctx context.Context, arg repository.CreateWorkoutAndReturnIdParams) (string, error)
	GetAll(ctx context.Context, userId string) ([]Workout, error)
	GetById(ctx context.Context, arg repository.GetWorkoutByIdParams) (Workout, error)
	DeleteById(ctx context.Context, arg repository.DeleteWorkoutByIdParams) error
	UpdateById(context context.Context, arg repository.UpdateWorkoutByIdParams) error
}

type workoutsRepository struct {
	repo repository.Querier
}

func (w *workoutsRepository) UpdateById(ctx context.Context, arg repository.UpdateWorkoutByIdParams) error {
	rows, err := w.repo.UpdateWorkoutById(ctx, arg)
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("workout not found")
	}
	return nil
}

func (w *workoutsRepository) DeleteById(ctx context.Context, arg repository.DeleteWorkoutByIdParams) error {
	rows, err := w.repo.DeleteWorkoutById(ctx, arg)
	if err != nil {
		return err
	}

	if rows == 0 {
		slog.Info("Tried to delete workout that did not exist", "workoutId", arg.ID)
	}
	return nil
}

func (w *workoutsRepository) CompleteById(ctx context.Context, arg repository.CompleteWorkoutByIdParams) (int64, error) {
	return w.repo.CompleteWorkoutById(ctx, arg)
}

func (w *workoutsRepository) CreateAndReturnId(ctx context.Context, arg repository.CreateWorkoutAndReturnIdParams) (string, error) {
	return w.repo.CreateWorkoutAndReturnId(ctx, arg)
}

func (w *workoutsRepository) GetAll(ctx context.Context, userId string) ([]Workout, error) {
	workouts, err := w.repo.GetAllWorkouts(ctx, userId)
	if err != nil {
		return []Workout{}, err
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
		return Workout{}, err
	}

	result := newWorkout(workout)
	return result, nil
}
