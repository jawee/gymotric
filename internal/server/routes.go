package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.helloWorldHandler)

	mux.HandleFunc("/health", s.healthHandler)

	mux.Handle("GET /workouts", http.HandlerFunc(s.getAllWorkoutsHandler))
	mux.Handle("GET /workouts/{id}", http.HandlerFunc(s.getWorkoutByIdHandler))
	mux.Handle("GET /workouts/{id}/exercises", http.HandlerFunc(s.getExercisesByWorkoutIdHandler))
	mux.Handle("GET /workouts/{id}/exercises/{exerciseId}/sets", http.HandlerFunc(s.getSetsByExerciseIdHandler))

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
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

func (s *Server) helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
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
