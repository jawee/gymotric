package exercisetypes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"weight-tracker/internal/database"
)

func AddEndpoints(mux *http.ServeMux, s database.Service, authenticationWrapper func(next http.Handler) http.Handler) {
	handler := handler{
		service: NewService(exerciseTypeRepository{s.GetRepository()}),
	}

	mux.Handle("GET /exercise-types", authenticationWrapper(http.HandlerFunc(handler.getAllWorkoutTypesHandler)))
	mux.Handle("POST /exercise-types", authenticationWrapper(http.HandlerFunc(handler.createExerciseTypeHandler)))
	mux.Handle("DELETE /exercise-types/{id}", authenticationWrapper(http.HandlerFunc(handler.deleteExerciseTypeByIdHandler)))
	mux.Handle("GET /exercise-types/{id}/max", authenticationWrapper(http.HandlerFunc(handler.getMaxSet)))
	mux.Handle("GET /exercise-types/{id}/last", authenticationWrapper(http.HandlerFunc(handler.getLastSet)))
}

type handler struct {
	service Service
}

func (s *handler) getLastSet(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	exerciseTypeId := r.PathValue("id")

	lastSet, err := s.service.GetLastWeightRepsByExerciseTypeId(r.Context(), exerciseTypeId, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		slog.Warn("Failed to get last set", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{"weight": lastSet.Weight, "reps": lastSet.Reps}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Warn("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
	}
}

func (s *handler) getMaxSet(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	exerciseTypeId := r.PathValue("id")

	lastSet, err := s.service.GetMaxWeightRepsByExerciseTypeId(r.Context(), exerciseTypeId, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		slog.Warn("Failed to get last set", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{"weight": lastSet.Weight, "reps": lastSet.Reps}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Warn("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
	}
}

func (s *handler) getAllWorkoutTypesHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	exerciseTypes, err := s.service.GetAll(r.Context(), userId)

	slog.Debug(fmt.Sprintf("returning %d exercise types", len(exerciseTypes)))

	resp := map[string]interface{}{"exercise_types": exerciseTypes}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Warn("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
	}
}

func (s *handler) deleteExerciseTypeByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	exerciseTypeId := r.PathValue("id")
	err := s.service.DeleteById(r.Context(), exerciseTypeId, userId)

	if err != nil {
		slog.Warn("Failed to delete exercise type", "error", err, "exerciseTypeId", exerciseTypeId)
		http.Error(w, "Failed to delete exercise type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (s *handler) createExerciseTypeHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	decoder := json.NewDecoder(r.Body)
	var t createExerciseTypeRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := s.service.CreateAndReturnId(r.Context(), t, userId)

	if err != nil {
		slog.Warn("Failed to create exercise", "error", err)
		http.Error(w, "Failed to create exercise", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	resp := map[string]interface{}{"id": id}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Warn("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
	}
}

type createExerciseTypeRequest struct {
	Name string `json:"name"`
}
