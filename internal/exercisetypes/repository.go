package exercisetypes

import (
	"context"
	"log/slog"
	"weight-tracker/internal/repository"
)

type ExerciseType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ExerciseTypeRepository interface {
	CreateAndReturnId(context context.Context, exerciseType repository.CreateExerciseTypeAndReturnIdParams) (string, error)
	DeleteById(context context.Context, exerciseTypeId string) error
	GetAll(context context.Context) ([]ExerciseType, error)
}

func (e exerciseTypeRepository) GetAll(context context.Context) ([]ExerciseType, error) {
	exerciseTypes, err := e.repo.GetAllExerciseTypes(context)
	if err != nil {
		return []ExerciseType{}, err
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

func (e exerciseTypeRepository) DeleteById(context context.Context, exerciseTypeId string) error {
	rows, err := e.repo.DeleteExerciseTypeById(context, exerciseTypeId)
	if err != nil {
		return err
	}

	if rows == 0 {
		slog.Info("Tried to delete exercise type that did not exist", "exerciseTypeId", exerciseTypeId)
	}
	return nil
}

func (e exerciseTypeRepository) CreateAndReturnId(context context.Context, exerciseType repository.CreateExerciseTypeAndReturnIdParams) (string, error) {
	return e.repo.CreateExerciseTypeAndReturnId(context, exerciseType)
}

type exerciseTypeRepository struct {
	repo repository.Querier
}
