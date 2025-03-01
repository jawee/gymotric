package users

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"weight-tracker/internal/database"
)

type createUserAndReturnIdRequest struct {
	Username string
	Password string
}

type loginRequest struct {
	Username string
	Password string
}

type handler struct {
	service Service
}

func AddEndpoints(mux *http.ServeMux, s database.Service) {
	handler := handler{
		service: NewService(&usersRepository{s.GetRepository()}),
	}
	// mux.Handle("GET /workouts/{id}/exercises/{exerciseId}/sets", authenticationWrapper(http.HandlerFunc(handler.getSetsByExerciseIdHandler)))
	// mux.Handle("POST /workouts/{id}/exercises/{exerciseId}/sets", authenticationWrapper(http.HandlerFunc(handler.createSetHandler)))
	// mux.Handle("DELETE /workouts/{id}/exercises/{exerciseId}/sets/{setId}", authenticationWrapper(http.HandlerFunc(handler.deleteSetByIdHandler)))

	mux.Handle("POST /users", http.HandlerFunc(handler.createUserHandler))
	mux.Handle("POST /login", http.HandlerFunc(handler.loginHandler))
}

func (s *handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t loginRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := s.service.Login(r.Context(), t)

	if err != nil {
		slog.Warn("Failed to login", "error", err)
		http.Error(w, "Failed to login", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{"token": token}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	cookie := http.Cookie{
        Name:     "weight-tracker-auth",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    }
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
	}
}

func (s *handler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t createUserAndReturnIdRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Warn("Invalid request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := s.service.CreateAndReturnId(r.Context(), t)

	if err != nil {
		slog.Warn("Failed to create user", "error", err)
		http.Error(w, "Failed to create user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	resp := map[string]interface{}{"id": id}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.Warn("Failed to write response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
	}
}

