package exercisetypes

import (
	"context"
	"sort"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

func NewService(repo ExerciseTypeRepository) Service {
	return &exerciseTypeService{repo}
}

type Service interface {
	GetAll(context context.Context) ([]ExerciseType, error)
	DeleteById(context context.Context, exerciseTypeId string) error
	CreateAndReturnId(context context.Context, exerciseType createExerciseTypeRequest) (string, error)
}

func (s *exerciseTypeService) CreateAndReturnId(context context.Context, exerciseType createExerciseTypeRequest) (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	toCreate := repository.CreateExerciseTypeAndReturnIdParams{
		ID:   uuid.String(),
		Name: exerciseType.Name,
	}

	id, err := s.repo.CreateAndReturnId(context, toCreate)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *exerciseTypeService) DeleteById(context context.Context, exerciseTypeId string) error {
	return s.repo.DeleteById(context, exerciseTypeId)
}

func (s *exerciseTypeService) GetAll(context context.Context) ([]ExerciseType, error) {
	exerciseTypes, err := s.repo.GetAll(context)
	if err != nil {
		return []ExerciseType{}, err
	}

	sort.Slice(exerciseTypes, func(i, j int) bool {
		return exerciseTypes[i].Name < exerciseTypes[j].Name
	})

	return exerciseTypes, nil
}

type exerciseTypeService struct {
	repo ExerciseTypeRepository
}
