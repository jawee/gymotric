package exerciseitems

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"weight-tracker/internal/database"
	"weight-tracker/internal/exercises"
	"weight-tracker/internal/repository"
	"weight-tracker/internal/utils"
)

func AddEndpoints(mux *http.ServeMux, s database.Service, authenticationWrapper func(next http.Handler) http.Handler) {
	handler := handler{
		service: NewService(
			NewExerciseItemRepository(s.GetRepository()),
			exercises.NewExerciseRepository(s.GetRepository()),
		),
	}

	mux.Handle("GET /workouts/{workoutId}/exercise-items", authenticationWrapper(http.HandlerFunc(handler.getByWorkoutIdHandler)))
	mux.Handle("GET /exercise-items/{id}", authenticationWrapper(http.HandlerFunc(handler.getByIdHandler)))
	mux.Handle("POST /workouts/{workoutId}/exercise-items", authenticationWrapper(http.HandlerFunc(handler.createHandler)))
	mux.Handle("PUT /exercise-items/{id}", authenticationWrapper(http.HandlerFunc(handler.updateHandler)))
	mux.Handle("DELETE /exercise-items/{id}", authenticationWrapper(http.HandlerFunc(handler.deleteHandler)))
}

type handler struct {
	service Service
}

type createExerciseItemRequest struct {
	Type string `json:"type"`
}

type updateExerciseItemRequest struct {
	Type string `json:"type"`
}

func (h *handler) getByWorkoutIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	workoutId := r.PathValue("workoutId")

	arg := repository.GetExerciseItemsByWorkoutIdParams{
		WorkoutID: workoutId,
		UserID:    userId,
	}
	items, err := h.service.GetByWorkoutIdWithExercises(r.Context(), arg)
	if err != nil {
		slog.Error("Failed to get exercise items by workout", "error", err, "workoutId", workoutId)
		http.Error(w, "Failed to get exercise items", http.StatusBadRequest)
		return
	}

	jsonResp, err := utils.CreateResponse(items)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	utils.ReturnJson(w, jsonResp)
}

func (h *handler) getByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	itemId := r.PathValue("id")

	arg := repository.GetExerciseItemByIdParams{
		ID:     itemId,
		UserID: userId,
	}
	item, err := h.service.GetByIdWithExercises(r.Context(), arg)
	if err != nil {
		slog.Error("Failed to get exercise item", "error", err, "itemId", itemId)
		http.Error(w, "Failed to get exercise item", http.StatusBadRequest)
		return
	}

	jsonResp, err := utils.CreateResponse(item)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	utils.ReturnJson(w, jsonResp)
}

func (h *handler) createHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	workoutId := r.PathValue("workoutId")

	var req createExerciseItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateAndReturnId(r.Context(), req.Type, workoutId, userId)
	if err != nil {
		slog.Error("Failed to create exercise item", "error", err)
		http.Error(w, "Failed to create exercise item", http.StatusBadRequest)
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

func (h *handler) updateHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	itemId := r.PathValue("id")

	var req updateExerciseItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	arg := repository.UpdateExerciseItemTypeParams{
		ID:        itemId,
		Type:      req.Type,
		UserID:    userId,
		UpdatedOn: timeNow(),
	}
	_, err := h.service.UpdateType(r.Context(), arg)
	if err != nil {
		slog.Error("Failed to update exercise item", "error", err)
		http.Error(w, "Failed to update exercise item", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)
	itemId := r.PathValue("id")

	arg := repository.DeleteExerciseItemByIdParams{
		ID:     itemId,
		UserID: userId,
	}
	_, err := h.service.DeleteById(r.Context(), arg)
	if err != nil {
		slog.Error("Failed to delete exercise item", "error", err)
		http.Error(w, "Failed to delete exercise item", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func timeNow() string {
	return time.Now().UTC().Format(time.RFC3339)
}
