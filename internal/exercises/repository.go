package exercises

import (
	"context"
	"log/slog"
	"weight-tracker/internal/exercisetypes"
	"weight-tracker/internal/repository"
)

type Exercise struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	WorkoutID      string `json:"workout_id"`
	ExerciseTypeID string `json:"exercise_type_id"`
}

type ExerciseRepository interface {
	GetAll(context context.Context) ([]Exercise, error)
	GetByWorkoutId(context context.Context, workoutId string) ([]Exercise, error)
	DeleteById(context context.Context, id string) error
	CreateAndReturnId(context context.Context, exercise repository.CreateExerciseAndReturnIdParams, workoutId string) (string, error)
	GetExerciseTypeById(context context.Context, exerciseTypeId string) (*exercisetypes.ExerciseType, error)
}

type exerciseRepository struct {
	repo repository.Querier
}

func (e exerciseRepository) CreateAndReturnId(context context.Context, exercise repository.CreateExerciseAndReturnIdParams, workoutId string) (string, error) {
	id, err := e.repo.CreateExerciseAndReturnId(context, exercise)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (e exerciseRepository) GetExerciseTypeById(context context.Context, exerciseTypeId string) (*exercisetypes.ExerciseType, error) {
	exerciseType, err := e.repo.GetExerciseTypeById(context, exerciseTypeId)
	if err != nil {
		slog.Warn("Failed GetExerciseTypeById", "error", err)
		return nil, err
	}

	return &exercisetypes.ExerciseType{ID: exerciseType.ID, Name: exerciseType.Name}, nil
}

func (e exerciseRepository) DeleteById(context context.Context, id string) error {
	rows, err := e.repo.DeleteExerciseById(context, id)

	if err != nil {
		return err
	}

	if rows == 0 {
		slog.Info("Tried to delete exercise that did not exist", "exerciseId", id)
	}

	return nil
}

func (e exerciseRepository) GetAll(context context.Context) ([]Exercise, error) {
	exercises, err := e.repo.GetAllExercises(context)

	if err != nil {
		return []Exercise{}, err
	}

	result := []Exercise{}
	for _, v := range exercises {
		result = append(result, newExercise(v))
	}

	return result, nil
}

func newExercise(v repository.Exercise) Exercise {
	exercise := Exercise{
		ID:             v.ID,
		ExerciseTypeID: v.ExerciseTypeID,
		Name:           v.Name,
		WorkoutID:      v.WorkoutID,
	}

	return exercise
}

func (e exerciseRepository) GetByWorkoutId(context context.Context, exerciseId string) ([]Exercise, error) {
	exercises, err := e.repo.GetExercisesByWorkoutId(context, exerciseId)
	slog.Debug("GetExercisesByWorkoutId returns", "exercises", exercises)

	if err != nil {
		return []Exercise{}, err
	}

	result := []Exercise{}
	for _, v := range exercises {
		result = append(result, newExercise(v))
	}

	slog.Debug("GetByWorkoutId returns", "exercises", result)
	return result, nil
}
