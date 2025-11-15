package exercises

import (
	"context"
	"fmt"
	"time"
	"weight-tracker/internal/exerciseitems"
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
	repo                   ExerciseRepository
	exerciseItemsService   exerciseitems.Service
}

func (e *exerciseService) CreateAndReturnId(context context.Context, exercise createExerciseRequest, workoutId string, userId string) (string, error) {
	arg := repository.GetExerciseTypeByIdParams{
		ID:     exercise.ExerciseTypeID,
		UserID: userId,
	}
	exerciseType, err := e.repo.GetExerciseTypeById(context, arg)

	if err != nil {
		return "", fmt.Errorf("failed to get exercise type by id: %w", err)
	}

	// Create exercise_item first
	exerciseItemId, err := e.exerciseItemsService.CreateAndReturnId(context, "exercise", workoutId, userId)
	if err != nil {
		return "", fmt.Errorf("failed to create exercise item: %w", err)
	}

	// Create exercise
	exerciseUUID, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID for exercise: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	toCreate := repository.CreateExerciseAndReturnIdParams{
		ID:             exerciseUUID.String(),
		Name:           exerciseType.Name,
		WorkoutID:      workoutId,
		ExerciseTypeID: exerciseType.ID,
		ExerciseItemID: exerciseItemId,
		CreatedOn:      now,
		UpdatedOn:      now,
		UserID:         userId,
	}

	id, err := e.repo.CreateAndReturnId(context, toCreate)
	if err != nil {
		return "", fmt.Errorf("failed to create exercise: %w", err)
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

func NewService(repo ExerciseRepository, exerciseItemsService exerciseitems.Service) Service {
	return &exerciseService{repo, exerciseItemsService}
}
