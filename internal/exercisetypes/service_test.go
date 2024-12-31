package exercisetypes

import (
	"context"
	"testing"
	"weight-tracker/internal/repository"
)

type repoMock struct {
	ExerciseTypes []ExerciseType
}

func (r *repoMock) CreateAndReturnId(context context.Context, exerciseType repository.CreateExerciseTypeAndReturnIdParams) (string, error) {
	exerciseTypeStruct := ExerciseType{
		ID: exerciseType.ID,
		Name: exerciseType.Name,
	}
	r.ExerciseTypes = append(r.ExerciseTypes, exerciseTypeStruct)

	return exerciseType.ID, nil
}

func (r *repoMock) DeleteById(context context.Context, exerciseTypeId string) error {
	for i, v := range r.ExerciseTypes {
		if v.ID == exerciseTypeId {
			r.ExerciseTypes = append(r.ExerciseTypes[:i], r.ExerciseTypes[i+1:]...)
		}
	}

	return nil
}

func (r *repoMock) GetAll(context context.Context) ([]ExerciseType, error) {
	return r.ExerciseTypes, nil
}

func TestGetAll(t *testing.T) {
	repoMock := repoMock{[]ExerciseType{
		{ ID: "b", Name: "b" },
		{ ID: "a", Name: "a" },
	}}
	service := NewService(&repoMock)

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

func TestDeleteById(t *testing.T) {
	repoMock := repoMock{[]ExerciseType{
		{ ID: "b", Name: "b" },
		{ ID: "a", Name: "a" },
	}}

	service := NewService(&repoMock)
	err := service.DeleteById(context.Background(), "a")

	if err != nil {
		t.Fatalf("Got err %s\n", err)
	}

	if len(repoMock.ExerciseTypes) != 1 {
		t.Fatalf("Expected repository to only have 1 exercise type, but has %d\n", len(repoMock.ExerciseTypes))
	}
}

func TestCreateAndReturnId(t *testing.T) {
	repoMock := repoMock{[]ExerciseType{
	}}

	service := NewService(&repoMock)
	_, err := service.CreateAndReturnId(context.Background(), createExerciseTypeRequest{
		Name: "a",
	})

	if err != nil {
		t.Fatalf("Got err %s\n", err)
	}

	if len(repoMock.ExerciseTypes) != 1 {
		t.Fatalf("Expected repository to only have 1 exercise type, but has %d\n", len(repoMock.ExerciseTypes))
	}
}
