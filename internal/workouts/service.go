package workouts

import (
	"context"
	"sort"
	"time"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	GetAll(context context.Context) ([]Workout, error)
	GetById(context context.Context, id string) (Workout, error)
	CreateAndReturnId(context context.Context, t createWorkoutRequest) (string, error)
	CompleteById(context context.Context, workoutId string) error
}

type workoutsService struct {
	repo WorkoutsRepository
}

func (w *workoutsService) CompleteById(context context.Context, workoutId string) error {
	completeParams := repository.CompleteWorkoutByIdParams{
		ID:          workoutId,
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
	}

	_, err := w.repo.CompleteById(context, completeParams)

	return err
}

func (w *workoutsService) CreateAndReturnId(context context.Context, t createWorkoutRequest) (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	workout := repository.CreateWorkoutAndReturnIdParams{
		ID:        uuid.String(),
		Name:      t.Name,
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
	}

	id, err := w.repo.CreateAndReturnId(context, workout)

	return id, err
}

func (w *workoutsService) GetAll(context context.Context) ([]Workout, error) {
	workouts, err := w.repo.GetAll(context)

	if err != nil {
		return []Workout{}, err
	}

	sort.Slice(workouts, func(i, j int) bool {
		ta, err := time.Parse(time.RFC3339, workouts[i].CreatedOn)

		if err != nil {
			return false
		}

		tb, err := time.Parse(time.RFC3339, workouts[j].CreatedOn)
		if err != nil {
			return false
		}

		return tb.Before(ta)
	})

	return workouts, nil
}

func (w *workoutsService) GetById(context context.Context, id string) (Workout, error) {
	workout, err := w.repo.GetById(context, id)
	return workout, err
}

func NewService(repo WorkoutsRepository) Service {
	return &workoutsService{repo}
}
