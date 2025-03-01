package sets

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

 func (r *repoMock) CreateAndReturnId(ctx context.Context, arg repository.CreateSetAndReturnIdParams) (string, error) {
 	args := r.Called(ctx, arg)
 	return args.String(0), args.Error(1)
 }

 func (r *repoMock) DeleteById(ctx context.Context, arg repository.DeleteSetByIdParams) (int64, error) {
 	args := r.Called(ctx, arg)
 	return args.Get(0).(int64), args.Error(1)
 }

 func (r *repoMock) GetAll(ctx context.Context, userId string) ([]Set, error) {
 	args := r.Called(ctx, userId)
 	return args.Get(0).([]Set), args.Error(1)
 }

 func (r *repoMock) GetById(ctx context.Context, arg repository.GetSetByIdParams) (Set, error) {
 	args := r.Called(ctx, arg)
 	return args.Get(0).(Set), args.Error(1)
 }

 func (r *repoMock) GetByExerciseId(ctx context.Context, arg repository.GetSetsByExerciseIdParams) ([]Set, error) {
 	args := r.Called(ctx, arg)
 	return args.Get(0).([]Set), args.Error(1)
 }

 func TestGetByExerciseId(t *testing.T) {
	userId := "userid"
	exerciseId := "exerciseIdA"
 	expected := []Set{
 		{ID: "a", Repetitions: 1, Weight: 10.0, ExerciseID: exerciseId},
 		{ID: "b", Repetitions: 1, Weight: 10.0, ExerciseID: exerciseId},
 	}

 	ctx := context.Background()

 	repoMock := repoMock{}
 	repoMock.On("GetByExerciseId", ctx, mock.MatchedBy(func(input repository.GetSetsByExerciseIdParams) bool {
		return input.ExerciseID == exerciseId && input.UserID == userId
	})).Return([]Set{
 		{ID: "a", Repetitions: 1, Weight: 10.0, ExerciseID: exerciseId},
 		{ID: "b", Repetitions: 1, Weight: 10.0, ExerciseID: exerciseId},
 	}, nil).Once()

 	service := NewService(&repoMock)

 	result, err := service.GetByExerciseId(ctx, exerciseId, userId)

 	assert.Nil(t, err)
 	assert.Len(t, result, 2)
 	assert.Equal(t, expected, result)
 	repoMock.AssertExpectations(t)
 }

 func TestCreateAndReturnId(t *testing.T) {
 	ctx := context.Background()

	userId := "userId"
 	setId := "setId"
 	exerciseId := "exerciseId"

 	repoMock := repoMock{}
 	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateSetAndReturnIdParams) bool {
 		return input.Weight == 10.5 && input.Repetitions == 1 && input.ExerciseID == exerciseId && input.CreatedOn != "" && input.UpdatedOn != "" && input.UserID == userId
 	})).Return(setId, nil).Once()

 	service := NewService(&repoMock)
 	id, err := service.CreateAndReturnId(context.Background(), createSetRequest{
 		Repetitions: 1,
 		Weight: 10.5,

 	}, exerciseId, userId)

 	assert.Nil(t, err)
 	assert.Equal(t, setId, id)
 	repoMock.AssertExpectations(t)
 }

func TestDeleteById(t *testing.T) {
	ctx := context.Background()
	userId := "userId"
	setId := "setId"

	repoMock := repoMock{}
	repoMock.On("DeleteById", ctx, mock.MatchedBy(func(input repository.DeleteSetByIdParams) bool {
		return input.ID == setId && input.UserID == userId
	})).Return(int64(1), nil).Once()

	service := NewService(&repoMock)
	err := service.DeleteById(ctx, setId, userId)

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}
