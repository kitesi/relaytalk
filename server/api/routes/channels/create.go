package channels

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	. "github.com/kitesi/relaytalk/api/routes/auth"
	"github.com/kitesi/relaytalk/db"
	. "github.com/kitesi/relaytalk/utils"
)

type CreateChannelRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateChannel(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		serverIDStr := chi.URLParam(r, "serverId")
		serverID, err := strconv.ParseInt(serverIDStr, 10, 32)

		if err != nil {
			SendJsonError(w, http.StatusBadRequest, "Invalid server ID")
		}

		var req CreateChannelRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = store.CreateChannel(r.Context(), db.CreateChannelParams{
			OwnerID:     int32(userID),
			Name:        req.Name,
			Description: ToPgText(req.Description),
			ServerID:    int32(serverID),
		})

		if err != nil {
			SendJsonError(w, http.StatusInternalServerError, "Failed to create channel")
			log.Printf("Error creating channel: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]string{
			"message": "Created channel successfully",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			SendJsonError(w, http.StatusInternalServerError, "Failed to encode response")
			return
		}
	}
}
