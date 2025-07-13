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
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userIdStr := strconv.Itoa(userID)

		w.Write([]byte("pong from protected endpoint, user ID: " + userIdStr))
	}
}
