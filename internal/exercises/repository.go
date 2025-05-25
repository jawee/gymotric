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
	GetAll(context context.Context, userId string) ([]Exercise, error)
	GetByWorkoutId(context context.Context, arg repository.GetExercisesByWorkoutIdParams) ([]Exercise, error)
	DeleteById(context context.Context, arg repository.DeleteExerciseByIdParams) error
	CreateAndReturnId(context context.Context, exercise repository.CreateExerciseAndReturnIdParams) (string, error)
	GetExerciseTypeById(context context.Context, arg repository.GetExerciseTypeByIdParams) (*exercisetypes.ExerciseType, error)
}

func NewExerciseRepository(repo repository.Querier) ExerciseRepository {
	return exerciseRepository{repo: repo}
}

type exerciseRepository struct {
	repo repository.Querier
}

func (e exerciseRepository) CreateAndReturnId(context context.Context, exercise repository.CreateExerciseAndReturnIdParams) (string, error) {
	id, err := e.repo.CreateExerciseAndReturnId(context, exercise)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (e exerciseRepository) GetExerciseTypeById(context context.Context, arg repository.GetExerciseTypeByIdParams) (*exercisetypes.ExerciseType, error) {
	exerciseType, err := e.repo.GetExerciseTypeById(context, arg)
	if err != nil {
		slog.Warn("Failed GetExerciseTypeById", "error", err)
		return nil, err
	}

	return &exercisetypes.ExerciseType{ID: exerciseType.ID, Name: exerciseType.Name}, nil
}

func (e exerciseRepository) DeleteById(context context.Context, arg repository.DeleteExerciseByIdParams) error {
	rows, err := e.repo.DeleteExerciseById(context, arg)

	if err != nil {
		return err
	}

	if rows == 0 {
		slog.Warn("Tried to delete exercise that did not exist", "exerciseId", arg.ID)
	}

	return nil
}

func (e exerciseRepository) GetAll(context context.Context, userId string) ([]Exercise, error) {
	exercises, err := e.repo.GetAllExercises(context, userId)

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

func (e exerciseRepository) GetByWorkoutId(context context.Context, arg repository.GetExercisesByWorkoutIdParams) ([]Exercise, error) {
	exercises, err := e.repo.GetExercisesByWorkoutId(context, arg)
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
