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

func (m *repoMock) DeleteById(context context.Context, arg repository.DeleteExerciseTypeByIdParams) error {
	args := m.Called(context, arg)
	return args.Error(0)
}

func (m *repoMock) GetAll(context context.Context, userId string) ([]ExerciseType, error) {
	args := m.Called(context, userId)
	return args.Get(0).([]ExerciseType), args.Error(1)
}

func (m *repoMock) GetLastWeightRepsByExerciseTypeId(ctx context.Context, arg repository.GetLastWeightRepsByExerciseTypeIdParams) (MaxLastWeightReps, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(MaxLastWeightReps), args.Error(1)
}

func (m *repoMock) GetMaxWeightRepsByExerciseTypeId(ctx context.Context, arg repository.GetMaxWeightRepsByExerciseTypeIdParams) (MaxLastWeightReps, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(MaxLastWeightReps), args.Error(1)
}

func (m *repoMock) UpdateById(ctx context.Context, arg repository.UpdateExerciseTypeParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func TestGetAll(t *testing.T) {
	userId := "userid"

	expected := []ExerciseType{
		{ID: "a", Name: "a"},
		{ID: "b", Name: "b"},
	}

	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetAll", ctx, userId).Return([]ExerciseType{
		{ID: "b", Name: "b"},
		{ID: "a", Name: "a"},
	}, nil).Once()

	service := NewService(&repoMock)

	result, err := service.GetAll(ctx, userId)

	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expected, result)
	repoMock.AssertExpectations(t)
}

func TestDeleteById(t *testing.T) {
	userId := "userid"
	exerciseTypeId := "exerciseTypeId"
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("DeleteById", ctx, mock.MatchedBy(func(input repository.DeleteExerciseTypeByIdParams) bool {
		return input.ID == exerciseTypeId && input.UserID == userId
	})).Return(nil).Once()

	service := NewService(&repoMock)
	err := service.DeleteById(ctx, exerciseTypeId, userId)

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestCreateAndReturnId(t *testing.T) {
	userId := "userid"
	exerciseTypeName := "exerciseTypeId"
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateExerciseTypeAndReturnIdParams) bool {
		return input.Name == exerciseTypeName && input.CreatedOn != "" && input.UpdatedOn != "" && input.UserID == userId
	})).Return("asdf", nil).Once()

	service := NewService(&repoMock)
	_, err := service.CreateAndReturnId(context.Background(), createExerciseTypeRequest{
		Name: exerciseTypeName,
	}, userId)

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestCreateTrimWhitespaceAndReturnId(t *testing.T) {
	userId := "userid"
	exerciseTypeName := " exerciseTypeId "
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateExerciseTypeAndReturnIdParams) bool {
		return input.Name == "exerciseTypeId" && input.CreatedOn != "" && input.UpdatedOn != "" && input.UserID == userId
	})).Return("asdf", nil).Once()

	service := NewService(&repoMock)
	_, err := service.CreateAndReturnId(context.Background(), createExerciseTypeRequest{
		Name: exerciseTypeName,
	}, userId)

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

