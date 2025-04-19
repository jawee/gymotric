package statistics

import (
	"database/sql"
	"log/slog"
	"net/http"
	"weight-tracker/internal/database"
	"weight-tracker/internal/utils"
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

type getStatisticsResponse struct {
	Week  int `json:"week"`
	Month int `json:"month"`
	Year  int `json:"year"`
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

	jsonResp, err := utils.CreateResponse(getStatisticsResponse{Week: statistics.Week, Month: statistics.Month, Year: statistics.Year})
	if err != nil {
		slog.Warn("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	utils.ReturnJson(w, jsonResp)
}
