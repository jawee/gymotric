package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"
	"weight-tracker/internal/exercises"
	"weight-tracker/internal/exercisetypes"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /health", http.HandlerFunc(s.healthHandler))

	exercisetypes.AddEndpoints(mux, s.db)

	mux.Handle("GET /workouts", http.HandlerFunc(s.getAllWorkoutsHandler))
	mux.Handle("POST /workouts", http.HandlerFunc(s.createWorkoutHandler))

	mux.Handle("GET /workouts/{id}", http.HandlerFunc(s.getWorkoutByIdHandler))
	mux.Handle("PUT /workouts/{id}/complete", http.HandlerFunc(s.completeWorkoutById))

	mux.Handle("GET /workouts/{id}/exercises/{exerciseId}/sets", http.HandlerFunc(s.getSetsByExerciseIdHandler))
	mux.Handle("POST /workouts/{id}/exercises/{exerciseId}/sets", http.HandlerFunc(s.createSetHandler))
	mux.Handle("DELETE /workouts/{id}/exercises/{exerciseId}/sets/{setId}", http.HandlerFunc(s.deleteSetByIdHandler))

	exercises.AddEndpoints(mux, s.db)

	return s.corsMiddleware(s.loggingMiddleware(mux))
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request", "Path", r.URL.Path, "Method", r.Method)

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) completeWorkoutById(w http.ResponseWriter, r *http.Request) {
	repo := s.db.GetRepository()

	workoutId := r.PathValue("id")

	completeParams := repository.CompleteWorkoutByIdParams{
		ID:          workoutId,
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
	}

	err := repo.CompleteWorkoutById(r.Context(), completeParams)

	if err != nil {
		slog.Warn("Failed to complete workout", "error", err, "workoutId", workoutId)
		http.Error(w, "Failed to complete workout", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
}


func (s *Server) deleteSetByIdHandler(w http.ResponseWriter, r *http.Request) {

	repo := s.db.GetRepository()

	setId := r.PathValue("setId")
	err := repo.DeleteSetById(r.Context(), setId)

	if err != nil {
		slog.Warn("Failed to delete set", "error", err, "setId", setId)
		http.Error(w, "Failed to delete set", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) createSetHandler(w http.ResponseWriter, r *http.Request) {
	repo := s.db.GetRepository()

	exerciseId := r.PathValue("exerciseId")
	decoder := json.NewDecoder(r.Body)
	var t createSetRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Warn("Failed to decode request body", "error", err)
		http.Error(w, "Failed to create workout", http.StatusBadRequest)
		return
	}

	set := repository.CreateSetAndReturnIdParams{
		ID:          generateUuid(),
		Repetitions: int64(t.Repetitions),
		Weight:      t.Weight,
		ExerciseID:  exerciseId,
	}
	id, err := repo.CreateSetAndReturnId(r.Context(), set)

	if err != nil {
		slog.Warn("Failed to create set", "error", err)
		http.Error(w, "Failed to create set", http.StatusBadRequest)
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

func (s *Server) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	repo := s.db.GetRepository()

	decoder := json.NewDecoder(r.Body)
	var t createWorkoutRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Warn("Failed to decode request body", "error", err)
		http.Error(w, "Failed to create workout", http.StatusBadRequest)
		return
	}

	workout := repository.CreateWorkoutAndReturnIdParams{
		ID:        generateUuid(),
		Name:      t.Name,
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
	}

	id, err := repo.CreateWorkoutAndReturnId(r.Context(), workout)

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

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		slog.Warn("Failed to write response", "error", err)
	}
}

func (s *Server) getAllWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	repo := s.db.GetRepository()
	workouts, err := repo.GetAllWorkouts(r.Context())

	sort.Slice(workouts, func(i, j int) bool {
		ta, err := time.Parse(time.RFC3339, workouts[i].CreatedOn)

		if err != nil {
			return false
		}

		tb, err := time.Parse(time.RFC3339, workouts[j].CreatedOn)
		if err != nil {
			return false
		}

		return tb.Before(ta)
	})

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

func (s *Server) getWorkoutByIdHandler(w http.ResponseWriter, r *http.Request) {
	repo := s.db.GetRepository()
	id := r.PathValue("id")
	workout, err := repo.GetWorkoutById(r.Context(), id)

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

func (s *Server) getSetsByExerciseIdHandler(w http.ResponseWriter, r *http.Request) {
	repo := s.db.GetRepository()
	id := r.PathValue("exerciseId")

	if id == "" {
		slog.Warn("Failed to get sets", "error", "empty exerciseId")
		http.Error(w, "Failed to get sets", http.StatusBadRequest)
		return
	}

	sets, err := repo.GetSetsByExerciseId(r.Context(), id)

	if err != nil {
		slog.Warn("Failed to get sets", "error", err)
		http.Error(w, "Failed to get sets", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{"sets": sets}
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

type createWorkoutRequest struct {
	Name string `json:"name"`
}

type createSetRequest struct {
	Repetitions int     `json:"repetitions"`
	Weight      float64 `json:"weight"`
}
