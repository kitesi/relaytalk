package utils

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

func SendJsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

func SendJsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

func ToPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{
			String: "",
			Valid:  false,
		}
	}
	return pgtype.Text{String: s, Valid: true}
}
