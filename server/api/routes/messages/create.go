package messages

import (
	"encoding/json"
	"net/http"

	. "github.com/kitesi/relaytalk/api/routes/auth"
	"github.com/kitesi/relaytalk/db"
	. "github.com/kitesi/relaytalk/utils"
)

type CreateMessageRequest struct {
	ChannelID int    `json:"channel_id"`
	Content   string `json:"content"`
}

func SendMessage(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req CreateMessageRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := store.CreateMessage(r.Context(), db.CreateMessageParams{
			UserID:    int32(userID),
			Content:   req.Content,
			ChannelID: int32(req.ChannelID),
		})

		if err != nil {
			SendJsonError(w, http.StatusInternalServerError, "Failed to send message")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]string{
			"message": "Message sent successfully",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			SendJsonError(w, http.StatusInternalServerError, "Failed to encode response")
			return
		}
	}
}
