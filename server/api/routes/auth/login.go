package auth

import (
	"encoding/json"
	"net/http"

	"github.com/kitesi/relaytalk/db"
	. "github.com/kitesi/relaytalk/utils"
	"golang.org/x/crypto/bcrypt"
)

func Login(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		json.NewDecoder(r.Body).Decode(&req)

		var userID int32
		var passwordHash string

		if req.Username == "" && req.Email == "" {
			SendJsonError(w, http.StatusBadRequest, "Username or email required")
			return
		}

		if req.Username != "" && req.Email != "" {
			SendJsonError(w, http.StatusBadRequest, "Please provide either username or email, not both")
			return
		}

		if req.Username != "" {
			user, err := store.GetUserByUsername(r.Context(), req.Username)

			if err != nil {
				SendJsonError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}

			userID = user.ID
			passwordHash = user.PasswordHash
		} else {
			user, err := store.GetUserByEmail(r.Context(), req.Email)

			if err != nil {
				SendJsonError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}

			userID = user.ID
			passwordHash = user.PasswordHash
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
			SendJsonError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		token, err := GenerateJWT(int(userID))

		if err != nil {
			SendJsonError(w, http.StatusInternalServerError, "Error generating token")
			return
		}

		SendJsonResponse(w, http.StatusOK, map[string]string{
			"token": token,
		})
	}
}
