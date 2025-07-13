package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kitesi/relaytalk/db"
	. "github.com/kitesi/relaytalk/utils"
	"golang.org/x/crypto/bcrypt"
)

func Register(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		json.NewDecoder(r.Body).Decode(&req)

		req.Username = strings.TrimSpace(req.Username)

		if req.Username == "" || req.Password == "" || req.Email == "" {
			SendJsonError(w, http.StatusBadRequest, "Username, password, and email are required")
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		err := store.CreateUser(r.Context(), db.CreateUserParams{
			Username:     req.Username,
			PasswordHash: string(hash),
			Email:        req.Email,
		})

		if err != nil {
			SendJsonError(w, http.StatusInternalServerError, "User already exists or database error")
			return
		}

		SendJsonResponse(w, http.StatusCreated, map[string]string{
			"message": "User registered successfully",
		})
	}
}
