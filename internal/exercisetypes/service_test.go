package exercisetypes

import (
	"context"
	"testing"
	"weight-tracker/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repoMock struct {
	mock.Mock
}

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
	assert.Equal(t, expected, result)
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
