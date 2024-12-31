package exercises

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"weight-tracker/internal/database"
)

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

	decoder := json.NewDecoder(r.Body)
	var t createExerciseRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Warn("Failed to decode request body", "error", err)
		http.Error(w, "Failed to create exercise", http.StatusBadRequest)
		return
	}

	workoutId := r.PathValue("id")

	id, err := s.service.CreateAndReturnId(r.Context(), t, workoutId)

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

	err := s.service.DeleteById(r.Context(), exerciseId)

	if err != nil {
		slog.Warn("Failed to delete exercise", "error", err, "exerciseId", exerciseId)
		http.Error(w, "Failed to delete exercise", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

type createExerciseRequest struct {
	ExerciseTypeID string `json:"exercise_type_id"`
}
