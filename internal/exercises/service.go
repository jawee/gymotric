package exercises

import (
	"context"
	"time"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	GetAll(context context.Context, userId string) ([]Exercise, error)
	GetByWorkoutId(context context.Context, workoutId string, userId string) ([]Exercise, error)
	DeleteById(context context.Context, id string, userId string) error
	CreateAndReturnId(context context.Context, exercise createExerciseRequest, workoutId string, userId string) (string, error)
}
type exerciseService struct {
	repo ExerciseRepository
}

func (e *exerciseService) CreateAndReturnId(context context.Context, exercise createExerciseRequest, workoutId string, userId string) (string, error) {
	arg := repository.GetExerciseTypeByIdParams{
		ID: exercise.ExerciseTypeID,
		UserID: userId,
	}
	exerciseType, err := e.repo.GetExerciseTypeById(context, arg)

	if err != nil {
		return "", err
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	toCreate := repository.CreateExerciseAndReturnIdParams{
		ID:             uuid.String(),
		Name:           exerciseType.Name,
		WorkoutID:      workoutId,
		ExerciseTypeID: exerciseType.ID,
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
		UserID: userId,
	}

	id, err := e.repo.CreateAndReturnId(context, toCreate)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (e *exerciseService) GetAll(context context.Context, userId string) ([]Exercise, error) {
	return e.repo.GetAll(context, userId)
}

func (e *exerciseService) GetByWorkoutId(context context.Context, workoutId string, userId string) ([]Exercise, error) {
	arg := repository.GetExercisesByWorkoutIdParams{
		WorkoutID: workoutId,
		UserID:    userId,
	}
	return e.repo.GetByWorkoutId(context, arg)
}

func (e *exerciseService) DeleteById(context context.Context, id string, userId string) error {
	arg := repository.DeleteExerciseByIdParams{
		ID:     id,
		UserID: userId,
	}
	return e.repo.DeleteById(context, arg)
}

func NewService(repo ExerciseRepository) Service {
	return &exerciseService{repo}
}
