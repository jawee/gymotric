package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func CreateResponse[T any](data T) ([]byte, error) {
	resp := map[string]T{
		"data": data,
	}
	return json.Marshal(resp)
}

func CreateIdResponse(id string) ([]byte, error) {
	resp := map[string]string{
		"id": id,
	}
	return json.Marshal(resp)
}

func CreatePaginatedResponse[T any](items []T, page int, pageSize int, totalCount int) ([]byte, error) {
	resp := map[string]any{
		"data":        items,
		"page":        page,
		"page_size":   pageSize,
		"total":       totalCount,
		"total_pages": (totalCount + pageSize - 1) / pageSize, // ceiling division
	}

	return json.Marshal(resp)
}

func CreateTokenResponse(token string, refreshToken string, tokenExpiration int) ([]byte, error) {
	resp := map[string]any{
		"access_token":  token,
		"token_type":    "Bearer",
		"expires_in":    tokenExpiration * 60,
		"refresh_token": refreshToken,
	}

	return json.Marshal(resp)
}

func ReturnJson(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		slog.Warn("Failed to write response", "error", err)
	}
}
