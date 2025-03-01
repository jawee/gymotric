package workouts

import (
	"context"
	"sort"
	"time"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	GetAll(context context.Context, userId string) ([]Workout, error)
	GetById(context context.Context, id string, userId string) (Workout, error)
	CreateAndReturnId(context context.Context, t createWorkoutRequest, userId string) (string, error)
	CompleteById(context context.Context, workoutId string, userId string) error
	DeleteById(context context.Context, workoutId string, userId string) error
}

type workoutsService struct {
	repo WorkoutsRepository
}

func (w *workoutsService) DeleteById(context context.Context, workoutId string, userId string) error {
	arg := repository.DeleteWorkoutByIdParams{
		ID:     workoutId,
		UserID: userId,
	}
	return w.repo.DeleteById(context, arg)
}

func (w *workoutsService) CompleteById(context context.Context, workoutId string, userId string) error {
	completeParams := repository.CompleteWorkoutByIdParams{
		ID:          workoutId,
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
		UserID:      userId,
	}

	_, err := w.repo.CompleteById(context, completeParams)

	return err
}

func (w *workoutsService) CreateAndReturnId(context context.Context, t createWorkoutRequest, userId string) (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	workout := repository.CreateWorkoutAndReturnIdParams{
		ID:        uuid.String(),
		Name:      t.Name,
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
		UserID:    userId,
	}

	id, err := w.repo.CreateAndReturnId(context, workout)

	return id, err
}

func (w *workoutsService) GetAll(context context.Context, userId string) ([]Workout, error) {
	workouts, err := w.repo.GetAll(context, userId)

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

func (w *workoutsService) GetById(context context.Context, id string, userId string) (Workout, error) {
	arg := repository.GetWorkoutByIdParams{
		ID:     id,
		UserID: userId,
	}
	workout, err := w.repo.GetById(context, arg)
	return workout, err
}

func NewService(repo WorkoutsRepository) Service {
	return &workoutsService{repo}
}
