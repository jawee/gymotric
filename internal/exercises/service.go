package exercises

import (
	"context"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	GetAll(context context.Context) ([]Exercise, error)
	GetByWorkoutId(context context.Context, workoutId string) ([]Exercise, error)
	DeleteById(context context.Context, id string) error
	CreateAndReturnId(context context.Context, exercise createExerciseRequest, workoutId string) (string, error)
}
type exerciseService struct {
	repo ExerciseRepository
}

func (e *exerciseService) CreateAndReturnId(context context.Context, exercise createExerciseRequest, workoutId string) (string, error) {
	exerciseType, err := e.repo.GetExerciseTypeById(context, exercise.ExerciseTypeID)

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
	}

	id, err := e.repo.CreateAndReturnId(context, toCreate, workoutId)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (e *exerciseService) GetAll(context context.Context) ([]Exercise, error) {
	return e.repo.GetAll(context)
}

func (e *exerciseService) GetByWorkoutId(context context.Context, workoutId string) ([]Exercise, error) {
	return e.repo.GetByWorkoutId(context, workoutId)
}

func (e *exerciseService) DeleteById(context context.Context, id string) error {
	return e.repo.DeleteById(context, id)
}

func NewService(repo ExerciseRepository) Service {
	return &exerciseService{repo}
}
