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

func AddEndpoints(mux *http.ServeMux, s database.Service, authenticationWrapper func(next http.Handler) http.Handler) {
	handler := handler{
		service: NewService(&workoutsRepository{s.GetRepository()}),
	}

	mux.Handle("GET /workouts", authenticationWrapper(http.HandlerFunc(handler.getAllWorkoutsHandler)))
	mux.Handle("POST /workouts", authenticationWrapper(http.HandlerFunc(handler.createWorkoutHandler)))
	mux.Handle("GET /workouts/{id}", authenticationWrapper(http.HandlerFunc(handler.getWorkoutByIdHandler)))
	mux.Handle("PUT /workouts/{id}/complete", authenticationWrapper(http.HandlerFunc(handler.completeWorkoutById)))
	mux.Handle("DELETE /workouts/{id}", authenticationWrapper(http.HandlerFunc(handler.deleteWorkoutByIdHandler)))
}

func (s *handler) deleteWorkoutByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	id := r.PathValue("id")
	err := s.service.DeleteById(r.Context(), id, userId)

	if err != nil {
		http.Error(w, "Failed to delete workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
}

func (s *handler) getAllWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	slog.Info("Getting all workouts")
	workouts, err := s.service.GetAll(r.Context(), userId)

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
	userId := r.Context().Value("sub").(string)
	id := r.PathValue("id")
	workout, err := s.service.GetById(r.Context(), id, userId)

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
	userId := r.Context().Value("sub").(string)
	decoder := json.NewDecoder(r.Body)
	var t createWorkoutRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Warn("Failed to decode request body", "error", err)
		http.Error(w, "Failed to create workout", http.StatusBadRequest)
		return
	}

	id, err := s.service.CreateAndReturnId(r.Context(), t, userId)

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
	userId := r.Context().Value("sub").(string)
	workoutId := r.PathValue("id")

	err := s.service.CompleteById(r.Context(), workoutId, userId)

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
