package exercises

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"weight-tracker/internal/database"
	"weight-tracker/internal/utils"
)

type handler struct {
	service Service
}

func AddEndpoints(mux *http.ServeMux, s database.Service, authenticationWrapper func(next http.Handler) http.Handler) {
	handler := handler{
		service: NewService(NewExerciseRepository(s.GetRepository())),
	}

	mux.Handle("GET /workouts/{workoutId}/exercise-items/{exerciseItemId}/exercises", authenticationWrapper(http.HandlerFunc(handler.getExercisesByWorkoutIdHandler)))
	mux.Handle("POST /workouts/{workoutId}/exercise-items/{exerciseItemId}/exercises", authenticationWrapper(http.HandlerFunc(handler.createExerciseHandler)))
	mux.Handle("DELETE /workouts/{workoutId}/exercise-items/{exerciseItemId}/exercises/{exerciseId}", authenticationWrapper(http.HandlerFunc(handler.deleteExerciseByIdHandler)))
}

func (s *handler) getExercisesByWorkoutIdHandler(w http.ResponseWriter, r *http.Request) {
	workoutId := r.PathValue("workoutId")

	userId := r.Context().Value("sub").(string)

	exercises, err := s.service.GetByWorkoutId(r.Context(), workoutId, userId)

	if err != nil {
		slog.Error("Failed to get exercises", "error", err)
		http.Error(w, "Failed to get exercises", http.StatusBadRequest)
		return
	}

	jsonResp, err := utils.CreateResponse(exercises)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	utils.ReturnJson(w, jsonResp)
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

	if t.ExerciseItemID == "" {
		slog.Error("Missing required field", "field", "exercise_item_id")
		http.Error(w, "exercise_item_id is required", http.StatusBadRequest)
		return
	}

	workoutId := r.PathValue("workoutId")

	id, err := s.service.CreateAndReturnId(r.Context(), t, workoutId, userId)

	if err != nil {
		slog.Error("Failed to create exercise", "error", err)
		http.Error(w, "Failed to create exercise", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	jsonResp, err := utils.CreateIdResponse(id)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	utils.ReturnJson(w, jsonResp)
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
	ExerciseItemID string `json:"exercise_item_id"`
}
