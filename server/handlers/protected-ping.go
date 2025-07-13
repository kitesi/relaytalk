package handlers

import (
	"net/http"
	"strconv"

	"github.com/kitesi/relaytalk/db"
)

func ProtectedPing(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())

		if !ok {
			sendJsonError(w, http.StatusUnauthorized, "Unauthorized access")
			return
		}

		userIdStr := strconv.Itoa(userID)

		w.Write([]byte("pong from protected endpoint, user ID: " + userIdStr))
	}
}
