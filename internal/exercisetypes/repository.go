package exercisetypes

import (
	"context"
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
	return e.repo.DeleteExerciseTypeById(context, exerciseTypeId)
}

func (e exerciseTypeRepository) CreateAndReturnId(context context.Context, exerciseType repository.CreateExerciseTypeAndReturnIdParams) (string, error) {
	return e.repo.CreateExerciseTypeAndReturnId(context, exerciseType)
}

type exerciseTypeRepository struct {
	repo repository.Querier
}
