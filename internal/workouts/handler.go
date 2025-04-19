package workouts

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"weight-tracker/internal/database"
	"weight-tracker/internal/exercises"
)

type handler struct {
	service Service
}

func AddEndpoints(mux *http.ServeMux, s database.Service, authenticationWrapper func(next http.Handler) http.Handler) {
	handler := handler{
		service: NewService(&workoutsRepository{s.GetRepository()}, exercises.NewExerciseRepository(s.GetRepository())),
	}

	mux.Handle("GET /workouts", authenticationWrapper(http.HandlerFunc(handler.getAllWorkoutsHandler)))
	mux.Handle("POST /workouts", authenticationWrapper(http.HandlerFunc(handler.createWorkoutHandler)))
	mux.Handle("GET /workouts/{id}", authenticationWrapper(http.HandlerFunc(handler.getWorkoutByIdHandler)))
	mux.Handle("PUT /workouts/{id}", authenticationWrapper(http.HandlerFunc(handler.updateWorkoutByIdHandler)))
	mux.Handle("PUT /workouts/{id}/complete", authenticationWrapper(http.HandlerFunc(handler.completeWorkoutById)))
	mux.Handle("POST /workouts/{id}/clone", authenticationWrapper(http.HandlerFunc(handler.cloneWorkoutById)))
	mux.Handle("DELETE /workouts/{id}", authenticationWrapper(http.HandlerFunc(handler.deleteWorkoutByIdHandler)))
}

func (s *handler) updateWorkoutByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	id := r.PathValue("id")
	decoder := json.NewDecoder(r.Body)
	var t updateWorkoutRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Failed to update workout", http.StatusBadRequest)
		return
	}

	slog.Debug("Updating workout", "id", id, "note", t.Note)

	err = s.service.UpdateWorkoutById(r.Context(), id, t, userId)

	if err != nil {
		slog.Error("Failed to update workout", "error", err)
		http.Error(w, "Failed to update workout", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *handler) cloneWorkoutById(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	id := r.PathValue("id")
	newId, err := s.service.CloneByIdAndReturnId(r.Context(), id, userId)
	if err != nil {
		slog.Error("Failed to create workout", "error", err)
		http.Error(w, "Failed to create workout", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	resp := map[string]any{"id": newId}
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

func (s *handler) deleteWorkoutByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	id := r.PathValue("id")
	err := s.service.DeleteById(r.Context(), id, userId)

	if err != nil {
		slog.Warn("Failed to delete workout", "error", err)
		http.Error(w, "Failed to delete workout", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
}

func createPaginatedResponse[T any](items []T, page int, pageSize int, totalCount int) ([]byte, error) {
	resp := map[string]any{
		"data":        items,
		"page":        page,
		"page_size":   pageSize,
		"total":       totalCount,
		"total_pages": (totalCount + pageSize - 1) / pageSize, // ceiling division
	}

	return json.Marshal(resp)
}

func createResponse[T any](data T) ([]byte, error) {
	resp := map[string]any{
		"data": data,
	}
	return json.Marshal(resp)
}

func createIdResponse(id string) ([]byte, error) {
	resp := map[string]any{
		"id": id,
	}
	return json.Marshal(resp)
}

func (s *handler) getAllWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	slog.Info("Getting all workouts")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // Default to 10 items per page
	}

	slog.Info("QueryParams", "page", page, "pageSize", pageSize)

	workouts, err := s.service.GetAll(r.Context(), userId, page, pageSize)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		slog.Error("Failed to get workout", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	count, err := s.service.GetAllCount(r.Context(), userId)
	if err != nil {
		slog.Error("Failed to get workout count", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	slog.Debug(fmt.Sprintf("returning %d workouts", len(workouts)))

	jsonResp, err := createPaginatedResponse(workouts, page, pageSize, count)

	if err != nil {
		slog.Warn("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	returnJson(w, jsonResp)
}

func returnJson(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		slog.Warn("Failed to write response", "error", err)
	}
}

var ErrNotFound = errors.New("sql: no rows in result set")

func (s *handler) getWorkoutByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	id := r.PathValue("id")
	workout, err := s.service.GetById(r.Context(), id, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		slog.Error("Failed to get workout", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	jsonResp, err := createResponse(workout)

	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	returnJson(w, jsonResp)
}

func (s *handler) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	decoder := json.NewDecoder(r.Body)
	var t createWorkoutRequest
	err := decoder.Decode(&t)

	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Failed to create workout", http.StatusBadRequest)
		return
	}

	id, err := s.service.CreateAndReturnId(r.Context(), t, userId)

	if err != nil {
		slog.Error("Failed to create workout", "error", err)
		http.Error(w, "Failed to create workout", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	jsonResp, err := createIdResponse(id)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	returnJson(w, jsonResp)
}

func (s *handler) completeWorkoutById(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	workoutId := r.PathValue("id")

	err := s.service.CompleteById(r.Context(), workoutId, userId)

	if err != nil {
		slog.Error("Failed to complete workout", "error", err, "workoutId", workoutId)
		http.Error(w, "Failed to complete workout", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
}

type createWorkoutRequest struct {
	Name string `json:"name"`
}

type updateWorkoutRequest struct {
	Note string `json:"note"`
}
