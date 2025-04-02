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
}

func (e exerciseTypeRepository) GetLastWeightRepsByExerciseTypeId(ctx context.Context, arg repository.GetLastWeightRepsByExerciseTypeIdParams) (MaxLastWeightReps, error) {
	a, err := e.repo.GetLastWeightRepsByExerciseTypeId(ctx, arg)
	if err != nil {
		return MaxLastWeightReps{}, err
	}

	return MaxLastWeightReps{ a.Weight, int(a.Repetitions) }, nil
}

func (e exerciseTypeRepository) GetMaxWeightRepsByExerciseTypeId(ctx context.Context, arg repository.GetMaxWeightRepsByExerciseTypeIdParams) (MaxLastWeightReps, error) {
	a, err := e.repo.GetMaxWeightRepsByExerciseTypeId(ctx, arg)
	if err != nil {
		return MaxLastWeightReps{}, err
	}

	weight := a.Weight
	repetitions := a.Repetitions.(int64)
	return MaxLastWeightReps{ weight, int(repetitions) }, nil
}

func (e exerciseTypeRepository) GetAll(context context.Context, userId string) ([]ExerciseType, error) {
	exerciseTypes, err := e.repo.GetAllExerciseTypes(context, userId)
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

func (e exerciseTypeRepository) DeleteById(context context.Context, arg repository.DeleteExerciseTypeByIdParams) error {
	rows, err := e.repo.DeleteExerciseTypeById(context, arg)
	if err != nil {
		return err
	}

	if rows == 0 {
		slog.Info("Tried to delete exercise type that did not exist", "exerciseTypeId", arg.ID)
	}
	return nil
}

func (e exerciseTypeRepository) CreateAndReturnId(context context.Context, exerciseType repository.CreateExerciseTypeAndReturnIdParams) (string, error) {
	return e.repo.CreateExerciseTypeAndReturnId(context, exerciseType)
}

type exerciseTypeRepository struct {
	repo repository.Querier
}
