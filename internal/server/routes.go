package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /health", http.HandlerFunc(s.healthHandler))

	mux.Handle("GET /exercise-types", http.HandlerFunc(s.getAllWorkoutTypesHandler))
	mux.Handle("POST /exercise-types", http.HandlerFunc(s.createExerciseTypeHandler))

	mux.Handle("GET /workouts", http.HandlerFunc(s.getAllWorkoutsHandler))
	mux.Handle("POST /workouts", http.HandlerFunc(s.createWorkoutHandler))

	mux.Handle("GET /workouts/{id}", http.HandlerFunc(s.getWorkoutByIdHandler))

	mux.Handle("GET /workouts/{id}/exercises", http.HandlerFunc(s.getExercisesByWorkoutIdHandler))
	mux.Handle("POST /workouts/{id}/exercises", http.HandlerFunc(s.createExerciseHandler))

	mux.Handle("GET /workouts/{id}/exercises/{exerciseId}/sets", http.HandlerFunc(s.getSetsByExerciseIdHandler))
	mux.Handle("POST /workouts/{id}/exercises/{exerciseId}/sets", http.HandlerFunc(s.createSetHandler))

	return s.corsMiddleware(s.loggingMiddleware(mux))
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
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

func (s *Server) getAllWorkoutTypesHandler(w http.ResponseWriter, r *http.Request) {
	repo := s.db.GetRepository()
	exerciseTypes, err := repo.GetAllExerciseTypes(r.Context())

	sort.Slice(exerciseTypes, func(i, j int) bool {
		return exerciseTypes[i].Name > exerciseTypes[j].Name
	})

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

func (s *Server) createSetHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) createExerciseHandler(w http.ResponseWriter, r *http.Request) {
	repo := s.db.GetRepository()

	decoder := json.NewDecoder(r.Body)
	var t createExerciseRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Warn("Failed to decode request body", "error", err)
		http.Error(w, "Failed to create exercise", http.StatusBadRequest)
		return
	}

	exerciseType, err := repo.GetExerciseTypeById(r.Context(), t.ExerciseTypeID)
	if err != nil {
		slog.Warn("Failed GetExerciseTypeById", "error", err)
		http.Error(w, "Failed to create exercise", http.StatusBadRequest)
		return
	}

	workoutId := r.PathValue("id")

	exercise := repository.CreateExerciseAndReturnIdParams{
		ID:             generateUuid(),
		Name:           exerciseType.Name,
		WorkoutID:      workoutId,
		ExerciseTypeID: exerciseType.ID,
	}

	id, err := repo.CreateExerciseAndReturnId(r.Context(), exercise)

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

func (s *Server) getExercisesByWorkoutIdHandler(w http.ResponseWriter, r *http.Request) {
	repo := s.db.GetRepository()
	id := r.PathValue("id")

	exercises, err := repo.GetExercisesByWorkoutId(r.Context(), id)

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

func (s *Server) getSetsByExerciseIdHandler(w http.ResponseWriter, r *http.Request) {
	repo := s.db.GetRepository()
	id := r.PathValue("exerciseId")

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

func (s *Server) createExerciseTypeHandler(w http.ResponseWriter, r *http.Request) {
	repo := s.db.GetRepository()

	decoder := json.NewDecoder(r.Body)
	var t createExerciseTypeRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Warn("Failed to decode request body", "error", err)
		http.Error(w, "Failed to create exercise type", http.StatusBadRequest)
		return
	}

	exerciseType := repository.CreateExerciseTypeAndReturnIdParams{
		ID:   generateUuid(),
		Name: t.Name,
	}

	id, err := repo.CreateExerciseTypeAndReturnId(r.Context(), exerciseType)

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

type createWorkoutRequest struct {
	Name string `json:"name"`
}

type createExerciseTypeRequest struct {
	Name string `json:"name"`
}

type createExerciseRequest struct {
	ExerciseTypeID string `json:"exercise_type_id"`
}
