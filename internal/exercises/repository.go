package exercises

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"weight-tracker/internal/database"
	"weight-tracker/internal/repository"
	"weight-tracker/internal/exercisetypes"

	"github.com/google/uuid"
)

type Exercise struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	WorkoutID      string `json:"workout_id"`
	ExerciseTypeID string `json:"exercise_type_id"`
}

type ExerciseRepository interface {
	GetAll(context context.Context) ([]Exercise, error)
	GetByWorkoutId(context context.Context, workoutId string) ([]Exercise, error)
	DeleteById(context context.Context, id string) error
	CreateAndReturnId(context context.Context, exercise repository.CreateExerciseAndReturnIdParams, workoutId string) (string, error)
	GetExerciseTypeById(context context.Context, exerciseTypeId string) (*exercisetypes.ExerciseType, error)
}

type exerciseRepository struct {
	repo repository.Querier
}

func (e exerciseRepository) CreateAndReturnId(context context.Context, exercise repository.CreateExerciseAndReturnIdParams, workoutId string) (string, error) {
	id, err := e.repo.CreateExerciseAndReturnId(context, exercise)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (e exerciseRepository) GetExerciseTypeById(context context.Context, exerciseTypeId string) (*exercisetypes.ExerciseType, error) {
	exerciseType, err := e.repo.GetExerciseTypeById(context, exerciseTypeId)
	if err != nil {
		slog.Warn("Failed GetExerciseTypeById", "error", err)
		return nil, err
	}

	return &exercisetypes.ExerciseType{ID: exerciseType.ID, Name: exerciseType.Name}, nil
}

func (e exerciseRepository) DeleteById(context context.Context, id string) error {
	rows, err := e.repo.DeleteExerciseById(context, id)

	if err != nil {
		return err
	}

	if rows == 0 {
		slog.Info("Tried to delete exercise that did not exist", "exerciseId", id)
	}

	return nil
}

func (e exerciseRepository) GetAll(context context.Context) ([]Exercise, error) {
	exercises, err := e.repo.GetAllExercises(context)

	if err != nil {
		return []Exercise{}, err
	}

	result := []Exercise{}
	for _, v := range exercises {
		result = append(result, newExercise(v))
	}

	return result, nil
}

func newExercise(v repository.Exercise) Exercise {
	exercise := Exercise{
		ID:             v.ID,
		ExerciseTypeID: v.ExerciseTypeID,
		Name:           v.Name,
		WorkoutID:      v.WorkoutID,
	}

	return exercise
}

func (e exerciseRepository) GetByWorkoutId(context context.Context, exerciseId string) ([]Exercise, error) {
	exercises, err := e.repo.GetExercisesByWorkoutId(context, exerciseId)
	slog.Debug("GetExercisesByWorkoutId returns", "exercises", exercises)

	if err != nil {
		return []Exercise{}, err
	}

	result := []Exercise{}
	for _, v := range exercises {
		result = append(result, newExercise(v))
	}

	slog.Debug("GetByWorkoutId returns", "exercises", result)
	return result, nil
}

type Service interface {
	GetAll(context context.Context) ([]Exercise, error)
	GetByWorkoutId(context context.Context, workoutId string) ([]Exercise, error)
	DeleteById(context context.Context, id string) error
	CreateAndReturnId(context context.Context, exercise createExerciseRequest, workoutId string) (string, error)
}
type exerciseService struct {
	repo ExerciseRepository
}

func (e *exerciseService) CreateAndReturnId(context context.Context, exercise createExerciseRequest, workoutId string) (string, error) {
	exerciseType, err := e.repo.GetExerciseTypeById(context, exercise.ExerciseTypeID)

	if err != nil {
		return "", err
	}

	toCreate := repository.CreateExerciseAndReturnIdParams{
		ID:             generateUuid(),
		Name:           exerciseType.Name,
		WorkoutID:      workoutId,
		ExerciseTypeID: exerciseType.ID,
	}

	id, err := e.repo.CreateAndReturnId(context, toCreate, workoutId)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (e *exerciseService) GetAll(context context.Context) ([]Exercise, error) {
	return e.repo.GetAll(context)
}

func (e *exerciseService) GetByWorkoutId(context context.Context, workoutId string) ([]Exercise, error) {
	return e.repo.GetByWorkoutId(context, workoutId)
}

func (e *exerciseService) DeleteById(context context.Context, id string) error {
	return e.repo.DeleteById(context, id)
}

func NewService(repo ExerciseRepository) Service {
	return &exerciseService{repo}
}

type handler struct {
	service Service
}

func AddEndpoints(mux *http.ServeMux, s database.Service) {
	handler := handler{
		service: NewService(exerciseRepository{s.GetRepository()}),
	}

	mux.Handle("GET /workouts/{id}/exercises", http.HandlerFunc(handler.getExercisesByWorkoutIdHandler))
	mux.Handle("POST /workouts/{id}/exercises", http.HandlerFunc(handler.createExerciseHandler))
	mux.Handle("DELETE /workouts/{id}/exercises/{exerciseId}", http.HandlerFunc(handler.deleteExerciseByIdHandler))
}

func (s *handler) getExercisesByWorkoutIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	exercises, err := s.service.GetByWorkoutId(r.Context(), id)

	if err != nil {
		slog.Warn("Failed to get exercises", "error", err)
		http.Error(w, "Failed to get exercises", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{"exercises": exercises}
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

func (s *handler) createExerciseHandler(w http.ResponseWriter, r *http.Request) {
	// repo := s.db.GetRepository()

	decoder := json.NewDecoder(r.Body)
	var t createExerciseRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Warn("Failed to decode request body", "error", err)
		http.Error(w, "Failed to create exercise", http.StatusBadRequest)
		return
	}

	// exerciseType, err := repo.GetExerciseTypeById(r.Context(), t.ExerciseTypeID)
	// if err != nil {
	// 	slog.Warn("Failed GetExerciseTypeById", "error", err)
	// 	http.Error(w, "Failed to create exercise", http.StatusBadRequest)
	// 	return
	// }

	workoutId := r.PathValue("id")

	// exercise := repository.CreateExerciseAndReturnIdParams{
	// 	ID:             generateUuid(),
	// 	Name:           exerciseType.Name,
	// 	WorkoutID:      workoutId,
	// 	ExerciseTypeID: exerciseType.ID,
	// }

	id, err := s.service.CreateAndReturnId(r.Context(), t, workoutId)
	// id, err := repo.CreateExerciseAndReturnId(r.Context(), exercise)

	if err != nil {
		slog.Warn("Failed to create exercise", "error", err)
		http.Error(w, "Failed to create exercise", http.StatusBadRequest)
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

func (s *handler) deleteExerciseByIdHandler(w http.ResponseWriter, r *http.Request) {
	exerciseId := r.PathValue("exerciseId")

	// err := repo.DeleteExerciseById(r.Context(), exerciseId)
	err := s.service.DeleteById(r.Context(), exerciseId)

	if err != nil {
		slog.Warn("Failed to delete exercise", "error", err, "exerciseId", exerciseId)
		http.Error(w, "Failed to delete exercise", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func generateUuid() string {
	id, _ := uuid.NewV7()
	return id.String()
}

type createExerciseRequest struct {
	ExerciseTypeID string `json:"exercise_type_id"`
}
