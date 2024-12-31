package workouts

import (
	"context"
	"weight-tracker/internal/repository"
)

type Workout struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	CompletedOn interface{} `json:"completed_on"`
	CreatedOn   string      `json:"created_on"`
	UpdatedOn   string      `json:"updated_on"`
}

type WorkoutsRepository interface {
	CompleteById(ctx context.Context, arg repository.CompleteWorkoutByIdParams) (int64, error)
	CreateAndReturnId(ctx context.Context, arg repository.CreateWorkoutAndReturnIdParams) (string, error)
	GetAll(ctx context.Context) ([]Workout, error)
	GetById(ctx context.Context, id string) (Workout, error)
}

type workoutsRepository struct {
	repo repository.Querier
}

func (w *workoutsRepository) CompleteById(ctx context.Context, arg repository.CompleteWorkoutByIdParams) (int64, error) {
	return w.repo.CompleteWorkoutById(ctx, arg)
}

func (w *workoutsRepository) CreateAndReturnId(ctx context.Context, arg repository.CreateWorkoutAndReturnIdParams) (string, error) {
	return w.repo.CreateWorkoutAndReturnId(ctx, arg)
}

func (w *workoutsRepository) GetAll(ctx context.Context) ([]Workout, error) {
	workouts, err := w.repo.GetAllWorkouts(ctx)
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
	return workout
}

func (w *workoutsRepository) GetById(ctx context.Context, id string) (Workout, error) {
	workout, err := w.repo.GetWorkoutById(ctx, id)
	if err != nil {
		return Workout{}, err
	}

	result := newWorkout(workout)
	return result, nil
}
