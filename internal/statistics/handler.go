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
	Week          int `json:"week"`
	PreviousWeek  int `json:"previous_week"`
	Month         int `json:"month"`
	PreviousMonth int `json:"previous_month"`
	Year          int `json:"year"`
	PreviousYear  int `json:"previous_year"`
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

	response := getStatisticsResponse{
		Week:          statistics.Week,
		PreviousWeek:  statistics.PreviousWeek,
		Month:         statistics.Month,
		PreviousMonth: statistics.PreviousMonth,
		Year:          statistics.Year,
		PreviousYear:  statistics.PreviousYear,
	}
	jsonResp, err := utils.CreateResponse(response)
	if err != nil {
		slog.Warn("Failed to marshal response", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	utils.ReturnJson(w, jsonResp)
}
