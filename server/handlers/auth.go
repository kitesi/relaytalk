package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kitesi/relaytalk/db"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

type contextKey string

const userContextKey = contextKey("userID")

func AuthMiddleware(store *db.Queries, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			sendJsonError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenStr)

		if err != nil {
			sendJsonError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userContextKey).(int)
	return userID, ok
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func Register(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		json.NewDecoder(r.Body).Decode(&req)

		req.Username = strings.TrimSpace(req.Username)

		if req.Username == "" || req.Password == "" || req.Email == "" {
			sendJsonError(w, http.StatusBadRequest, "Username, password, and email are required")
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		err := store.CreateUser(r.Context(), db.CreateUserParams{
			Username:     req.Username,
			PasswordHash: string(hash),
			Email:        req.Email,
		})

		if err != nil {
			sendJsonError(w, http.StatusInternalServerError, "User already exists or database error")
			return
		}

		sendJsonResponse(w, http.StatusCreated, map[string]string{
			"message": "User registered successfully",
		})
	}
}

func Login(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		json.NewDecoder(r.Body).Decode(&req)

		var userID int32
		var passwordHash string

		if req.Username == "" && req.Email == "" {
			sendJsonError(w, http.StatusBadRequest, "Username or email required")
			return
		}

		if req.Username != "" && req.Email != "" {
			sendJsonError(w, http.StatusBadRequest, "Please provide either username or email, not both")
			return
		}

		if req.Username != "" {
			user, err := store.GetUserByUsername(r.Context(), req.Username)

			if err != nil {
				sendJsonError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}

			userID = user.ID
			passwordHash = user.PasswordHash
		} else {
			user, err := store.GetUserByEmail(r.Context(), req.Email)

			if err != nil {
				sendJsonError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}

			userID = user.ID
			passwordHash = user.PasswordHash
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
			sendJsonError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		token, err := GenerateJWT(int(userID))

		if err != nil {
			sendJsonError(w, http.StatusInternalServerError, "Error generating token")
			return
		}

		sendJsonResponse(w, http.StatusOK, map[string]string{
			"token": token,
		})
	}
}
