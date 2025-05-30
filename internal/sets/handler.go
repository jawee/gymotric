package sets

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
		service: NewService(&setsRepository{s.GetRepository()}),
	}
	mux.Handle("GET /workouts/{id}/exercises/{exerciseId}/sets", authenticationWrapper(http.HandlerFunc(handler.getSetsByExerciseIdHandler)))
	mux.Handle("POST /workouts/{id}/exercises/{exerciseId}/sets", authenticationWrapper(http.HandlerFunc(handler.createSetHandler)))
	mux.Handle("DELETE /workouts/{id}/exercises/{exerciseId}/sets/{setId}", authenticationWrapper(http.HandlerFunc(handler.deleteSetByIdHandler)))
}

func (s *handler) deleteSetByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	setId := r.PathValue("setId")

	err := s.service.DeleteById(r.Context(), setId, userId)

	if err != nil {
		slog.Error("Failed to delete set", "error", err, "setId", setId)
		http.Error(w, "Failed to delete set", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (s *handler) createSetHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	exerciseId := r.PathValue("exerciseId")
	decoder := json.NewDecoder(r.Body)
	var t createSetRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := s.service.CreateAndReturnId(r.Context(), t, exerciseId, userId)

	if err != nil {
		slog.Warn("Failed to create set", "error", err)
		http.Error(w, "Failed to create set", http.StatusBadRequest)
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

func (s *handler) getSetsByExerciseIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	exerciseId := r.PathValue("exerciseId")

	if exerciseId == "" {
		slog.Warn("Failed to get sets", "error", "empty exerciseId")
		http.Error(w, "Failed to get sets", http.StatusBadRequest)
		return
	}

	sets, err := s.service.GetByExerciseId(r.Context(), exerciseId, userId)

	if err != nil {
		slog.Warn("Failed to get sets", "error", err)
		http.Error(w, "Failed to get sets", http.StatusBadRequest)
		return
	}

	jsonResp, err := utils.CreateResponse(sets)
	if err != nil {
		slog.Warn("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	utils.ReturnJson(w, jsonResp)
}

type createSetRequest struct {
	Repetitions int     `json:"repetitions"`
	Weight      float64 `json:"weight"`
}
