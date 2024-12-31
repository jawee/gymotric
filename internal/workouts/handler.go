package workouts

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"weight-tracker/internal/database"
)

type handler struct {
	service Service
}

func AddEndpoints(mux *http.ServeMux, s database.Service) {
	handler := handler{
		service: NewService(&workoutsRepository{s.GetRepository()}),
	}

	mux.Handle("GET /workouts", http.HandlerFunc(handler.getAllWorkoutsHandler))
	mux.Handle("POST /workouts", http.HandlerFunc(handler.createWorkoutHandler))
	mux.Handle("GET /workouts/{id}", http.HandlerFunc(handler.getWorkoutByIdHandler))
	mux.Handle("PUT /workouts/{id}/complete", http.HandlerFunc(handler.completeWorkoutById))
}

func (s *handler) getAllWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	workouts, err := s.service.GetAll(r.Context())

	slog.Debug(fmt.Sprintf("returning %d workouts", len(workouts)))

	resp := map[string]interface{}{"workouts": workouts}
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

func (s *handler) getWorkoutByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	workout, err := s.service.GetById(r.Context(), id)

	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{"workout": workout}
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

func (s *handler) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t createWorkoutRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Warn("Failed to decode request body", "error", err)
		http.Error(w, "Failed to create workout", http.StatusBadRequest)
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

func (s *handler) completeWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId := r.PathValue("id")

	err := s.service.CompleteById(r.Context(), workoutId)

	if err != nil {
		slog.Warn("Failed to complete workout", "error", err, "workoutId", workoutId)
		http.Error(w, "Failed to complete workout", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
}

type createWorkoutRequest struct {
	Name string `json:"name"`
}
