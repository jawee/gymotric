package exercisetypes

import (
	"context"
	"testing"
	"weight-tracker/internal/repository"
)

// type Service interface {
// 	GetAll(context context.Context) ([]ExerciseType, error)
// 	DeleteById(context context.Context, exerciseTypeId string) error
// 	CreateAndReturnId(context context.Context, exerciseType createExerciseTypeRequest) (string, error)
// }

type repoMock struct {
	exerciseTypes []ExerciseType
}

func (r repoMock) CreateAndReturnId(context context.Context, exerciseType repository.CreateExerciseTypeAndReturnIdParams) (string, error) {
	panic("unimplemented")
}

func (r repoMock) DeleteById(context context.Context, exerciseTypeId string) error {
	panic("unimplemented")
}

func (r repoMock) GetAll(context context.Context) ([]ExerciseType, error) {
	return r.exerciseTypes, nil
}

func TestGetAll(t *testing.T) {
	repoMock := repoMock{[]ExerciseType{
		{ ID: "b", Name: "b" },
		{ ID: "a", Name: "a" },
	}}
	service := NewService(repoMock)

	result, err := service.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Got err %s\n", err)
	}

	if len(result) != 2 {
		t.Fatalf("Got %d results, expected 2\n", len(result))
	}

	if result[0].Name != "a" {
		t.Fatalf("Not sorted\n")
	}
}
