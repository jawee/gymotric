package statistics

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"weight-tracker/internal/database"
)

func AddEndpoints(mux *http.ServeMux, s database.Service, authenticationWrapper func(next http.Handler) http.Handler) {
	handler := handler{
		service: NewService(&statisticsRepository{s.GetRepository()}),
	}

	mux.Handle("GET /statistics", authenticationWrapper(http.HandlerFunc(handler.getStatistics)))
}

type handler struct {
	service Service
}

func (s *handler) getStatistics(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("sub").(string)

	statistics, err := s.service.GetStatistics(r.Context(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		slog.Warn("Failed to get statistics", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resp := map[string]any{"week": statistics.Week, "month": statistics.Month, "year": statistics.Year}
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
