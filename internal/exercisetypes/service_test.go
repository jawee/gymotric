package exercisetypes

import (
	"context"
	"testing"
	"weight-tracker/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repoMock struct {
	// ExerciseTypes []ExerciseType
	mock.Mock
}

//
// func (r *repoMock) CreateAndReturnId(context context.Context, exerciseType repository.CreateExerciseTypeAndReturnIdParams) (string, error) {
// 	exerciseTypeStruct := ExerciseType{
// 		ID: exerciseType.ID,
// 		Name: exerciseType.Name,
// 	}
// 	r.ExerciseTypes = append(r.ExerciseTypes, exerciseTypeStruct)
//
// 	return exerciseType.ID, nil
// }
//
// func (r *repoMock) DeleteById(context context.Context, exerciseTypeId string) error {
// 	for i, v := range r.ExerciseTypes {
// 		if v.ID == exerciseTypeId {
// 			r.ExerciseTypes = append(r.ExerciseTypes[:i], r.ExerciseTypes[i+1:]...)
// 		}
// 	}
//
// 	return nil
// }
//
// func (r *repoMock) GetAll(context context.Context) ([]ExerciseType, error) {
// 	return r.ExerciseTypes, nil
// }

func (m *repoMock) CreateAndReturnId(context context.Context, exerciseType repository.CreateExerciseTypeAndReturnIdParams) (string, error) {
	args := m.Called(context, exerciseType)
	return args.String(0), args.Error(1)
}

func (m *repoMock) DeleteById(context context.Context, exerciseTypeId string) error {
	args := m.Called(context, exerciseTypeId)
	return args.Error(0)
}

func (m *repoMock) GetAll(context context.Context) ([]ExerciseType, error) {
	args := m.Called(context)
	return args.Get(0).([]ExerciseType), args.Error(1)
}

func TestGetAll(t *testing.T) {
	expected := []ExerciseType{
		{ID: "a", Name: "a"},
		{ID: "b", Name: "b"},
	}

	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetAll", ctx).Return([]ExerciseType{
		{ID: "b", Name: "b"},
		{ID: "a", Name: "a"},
	}, nil).Once()

	service := NewService(&repoMock)

	result, err := service.GetAll(ctx)

	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, result, expected)
	repoMock.AssertExpectations(t)
}

func TestDeleteById(t *testing.T) {
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("DeleteById", ctx, "a").Return(nil).Once()

	service := NewService(&repoMock)
	err := service.DeleteById(ctx, "a")

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestCreateAndReturnId(t *testing.T) {
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func (input repository.CreateExerciseTypeAndReturnIdParams) bool {
	return input.Name == "a" })).Return("asdf", nil).Once()

	service := NewService(&repoMock)
	_, err := service.CreateAndReturnId(context.Background(), createExerciseTypeRequest{
		Name: "a",
	})

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}
