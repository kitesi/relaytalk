package servers

import (
	"encoding/json"
	"net/http"

	. "github.com/kitesi/relaytalk/api/routes/auth"
	"github.com/kitesi/relaytalk/db"
	. "github.com/kitesi/relaytalk/utils"
)

type CreateServerRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateServer(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req CreateServerRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := store.CreateServer(r.Context(), db.CreateServerParams{
			OwnerID:     int32(userID),
			Name:        req.Name,
			Description: ToPgText(req.Description),
		})

		if err != nil {
			SendJsonError(w, http.StatusInternalServerError, "Failed to send message")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]string{
			"message": "Created server successfully",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			SendJsonError(w, http.StatusInternalServerError, "Failed to encode response")
			return
		}
	}
}
