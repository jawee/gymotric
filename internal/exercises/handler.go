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

func AddEndpoints(mux *http.ServeMux, s database.Service, authenticationWrapper func(next http.Handler) http.Handler) {
	handler := handler{
		service: NewService(NewExerciseRepository(s.GetRepository())),
	}

	mux.Handle("GET /workouts/{id}/exercises", authenticationWrapper(http.HandlerFunc(handler.getExercisesByWorkoutIdHandler)))
	mux.Handle("POST /workouts/{id}/exercises", authenticationWrapper(http.HandlerFunc(handler.createExerciseHandler)))
	mux.Handle("DELETE /workouts/{id}/exercises/{exerciseId}", authenticationWrapper(http.HandlerFunc(handler.deleteExerciseByIdHandler)))
}

func (s *handler) getExercisesByWorkoutIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	userId := r.Context().Value("sub").(string)

	exercises, err := s.service.GetByWorkoutId(r.Context(), id, userId)

	if err != nil {
		slog.Error("Failed to get exercises", "error", err)
		http.Error(w, "Failed to get exercises", http.StatusBadRequest)
		return
	}

	resp := map[string]any{"exercises": exercises}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
	}
}

func (s *handler) createExerciseHandler(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value("sub").(string)
	decoder := json.NewDecoder(r.Body)
	var t createExerciseRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workoutId := r.PathValue("id")

	id, err := s.service.CreateAndReturnId(r.Context(), t, workoutId, userId)

	if err != nil {
		slog.Error("Failed to create exercise", "error", err)
		http.Error(w, "Failed to create exercise", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	resp := map[string]any{"id": id}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
	}
}

func (s *handler) deleteExerciseByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	exerciseId := r.PathValue("exerciseId")

	err := s.service.DeleteById(r.Context(), exerciseId, userId)

	if err != nil {
		slog.Error("Failed to delete exercise", "error", err, "exerciseId", exerciseId)
		http.Error(w, "Failed to delete exercise", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

type createExerciseRequest struct {
	ExerciseTypeID string `json:"exercise_type_id"`
}
