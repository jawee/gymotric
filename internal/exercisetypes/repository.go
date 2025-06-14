package exercisetypes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"weight-tracker/internal/repository"
)

type ExerciseType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MaxLastWeightReps struct {
	Weight float64 `json:"weight"`
	Reps   int     `json:"reps"`
}

type ExerciseTypeRepository interface {
	CreateAndReturnId(context context.Context, exerciseType repository.CreateExerciseTypeAndReturnIdParams) (string, error)
	DeleteById(context context.Context, arg repository.DeleteExerciseTypeByIdParams) error
	GetAll(context context.Context, userId string) ([]ExerciseType, error)
	GetLastWeightRepsByExerciseTypeId(ctx context.Context, arg repository.GetLastWeightRepsByExerciseTypeIdParams) (MaxLastWeightReps, error)
	GetMaxWeightRepsByExerciseTypeId(ctx context.Context, arg repository.GetMaxWeightRepsByExerciseTypeIdParams) (MaxLastWeightReps, error)
	UpdateById(ctx context.Context, arg repository.UpdateExerciseTypeParams) error
}

func (e exerciseTypeRepository) UpdateById(ctx context.Context, arg repository.UpdateExerciseTypeParams) error {
	rows, err := e.repo.UpdateExerciseType(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to update exercise type: %w", err)
	}

	if rows == 0 {
		slog.Warn("Tried to update exercise type that did not exist", "exerciseTypeId", arg.ID)
		return errors.New("exercise type not found")
	}
	return nil
}

func (e exerciseTypeRepository) GetLastWeightRepsByExerciseTypeId(ctx context.Context, arg repository.GetLastWeightRepsByExerciseTypeIdParams) (MaxLastWeightReps, error) {
	a, err := e.repo.GetLastWeightRepsByExerciseTypeId(ctx, arg)
	if err != nil {
		return MaxLastWeightReps{}, fmt.Errorf("failed to get last weight reps by exercise type id: %w", err)
	}

	return MaxLastWeightReps{ a.Weight, int(a.Repetitions) }, nil
}

func (e exerciseTypeRepository) GetMaxWeightRepsByExerciseTypeId(ctx context.Context, arg repository.GetMaxWeightRepsByExerciseTypeIdParams) (MaxLastWeightReps, error) {
	a, err := e.repo.GetMaxWeightRepsByExerciseTypeId(ctx, arg)
	if err != nil {
		return MaxLastWeightReps{}, fmt.Errorf("failed to get max weight reps by exercise type id: %w", err)
	}

	weight := a.Weight
	repetitions := a.Repetitions.(int64)
	return MaxLastWeightReps{ weight, int(repetitions) }, nil
}

func (e exerciseTypeRepository) GetAll(context context.Context, userId string) ([]ExerciseType, error) {
	exerciseTypes, err := e.repo.GetAllExerciseTypes(context, userId)
	if err != nil {
		return []ExerciseType{}, fmt.Errorf("failed to get all exercise types: %w", err)
	}

	result := []ExerciseType{}
	for _, v := range exerciseTypes {
		result = append(result, newExerciseType(v))
	}
	return result, nil
}

func newExerciseType(v repository.ExerciseType) ExerciseType {
	return ExerciseType{
		ID:   v.ID,
		Name: v.Name,
	}
}

func (e exerciseTypeRepository) DeleteById(context context.Context, arg repository.DeleteExerciseTypeByIdParams) error {
	rows, err := e.repo.DeleteExerciseTypeById(context, arg)
	if err != nil {
		return fmt.Errorf("failed to delete exercise type: %w", err)
	}

	if rows == 0 {
		slog.Warn("Tried to delete exercise type that did not exist", "exerciseTypeId", arg.ID)
	}
	return nil
}

func (e exerciseTypeRepository) CreateAndReturnId(context context.Context, exerciseType repository.CreateExerciseTypeAndReturnIdParams) (string, error) {
	return e.repo.CreateExerciseTypeAndReturnId(context, exerciseType)
}

type exerciseTypeRepository struct {
	repo repository.Querier
}
