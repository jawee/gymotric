package exercisetypes

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"weight-tracker/internal/database"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
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

func AddEndpoints(mux *http.ServeMux, s database.Service) {
	handler := handler{
		service: NewService(exerciseTypeRepository{s.GetRepository()}),
	}

	mux.Handle("GET /exercise-types", http.HandlerFunc(handler.getAllWorkoutTypesHandler))
	mux.Handle("POST /exercise-types", http.HandlerFunc(handler.createExerciseTypeHandler))
	mux.Handle("DELETE /exercise-types/{id}", http.HandlerFunc(handler.deleteExerciseTypeByIdHandler))
}

func NewService(repo ExerciseTypeRepository) Service {
	return &exerciseTypeService{repo}
}

type Service interface {
	GetAll(context context.Context) ([]ExerciseType, error)
	DeleteById(context context.Context, exerciseTypeId string) error
	CreateAndReturnId(context context.Context, exerciseType createExerciseTypeRequest) (string, error)
}

func (s *exerciseTypeService) CreateAndReturnId(context context.Context, exerciseType createExerciseTypeRequest) (string, error) {
	toCreate := repository.CreateExerciseTypeAndReturnIdParams{
		ID:   generateUuid(),
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

type handler struct {
	service Service
}

func (s *handler) getAllWorkoutTypesHandler(w http.ResponseWriter, r *http.Request) {
	exerciseTypes, err := s.service.GetAll(r.Context())

	slog.Debug(fmt.Sprintf("returning %d exercise types", len(exerciseTypes)))

	resp := map[string]interface{}{"exercise_types": exerciseTypes}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
	}
}

func (s *handler) deleteExerciseTypeByIdHandler(w http.ResponseWriter, r *http.Request) {
	exerciseTypeId := r.PathValue("id")
	err := s.service.DeleteById(r.Context(), exerciseTypeId)

	if err != nil {
		slog.Warn("Failed to delete exercise type", "error", err, "exerciseTypeId", exerciseTypeId)
		http.Error(w, "Failed to delete exercise type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (s *handler) createExerciseTypeHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t createExerciseTypeRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Warn("Failed to decode request body", "error", err)
		http.Error(w, "Failed to create exercise type", http.StatusBadRequest)
		return
	}

	id, err := s.service.CreateAndReturnId(r.Context(), t)

	if err != nil {
		slog.Warn("Failed to create workout", "error", err)
		http.Error(w, "Failed to create workout", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	resp := map[string]interface{}{"id": id}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
	}
}

func generateUuid() string {
	id, _ := uuid.NewV7()
	return id.String()
}

type createExerciseTypeRequest struct {
	Name string `json:"name"`
}
