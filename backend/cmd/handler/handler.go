package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SpectreFury/odin-book/backend/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

type AuthHandler struct {
	DB *pgxpool.Pool
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	JwtSecret := []byte(os.Getenv("JWT_SECRET"))
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(user.FirstName)) < 1 || len(strings.TrimSpace(user.LastName)) < 1 || len(strings.TrimSpace(user.Email)) < 1 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Unable to hash password", err)
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
		return
	}

	var newID int32
	err = h.DB.QueryRow(
		context.Background(),
		`INSERT INTO users
		(first_name, last_name, email, password)
		VALUES
		($1, $2, $3, $4)
		RETURNING id`, user.FirstName, user.LastName, user.Email, hashedPassword).Scan(&newID)

	if err != nil {
		logger.Error("Duplicate id", err)
		http.Error(w, "You already have an account", http.StatusBadRequest)
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "odin-book",
		"sub": "userid",
		"id":  newID,
		"exp": time.Now().Add(time.Hour + 72).Unix(),
	})

	token, err := t.SignedString(JwtSecret)
	if err != nil {
		logger.Error("Unable to sign JWT", err)
		http.Error(w, "Unable to sign JWT", http.StatusInternalServerError)
		return
	}

	response := Response{
		Status:  true,
		Message: "Successfully signed up",
		Data: map[string]any{
			"token": token,
			"id":    newID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	JwtSecret := []byte(os.Getenv("JWT_SECRET"))

	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		logger.Error("Unable to decode json", err)
		http.Error(w, "Malformed body", http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(credentials.Email)) < 1 || len(strings.TrimSpace(credentials.Password)) < 1 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	var userID int64
	var hashedPassword string
	err = h.DB.QueryRow(context.Background(), `SELECT id, password FROM users WHERE email = $1`, credentials.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		logger.Error("No user found with email: ", credentials.Email, err)
		http.Error(w, "Unable to find user with the email", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
	if err != nil {
		logger.Error("Incorrect password", credentials.Email, err)
		http.Error(w, "Either email or password is incorrect", http.StatusUnauthorized)
		return
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "odin-book",
		"sub": "userid",
		"id":  userID,
		"exp": time.Now().Add(time.Hour + 72).Unix(),
	})

	token, err := t.SignedString(JwtSecret)
	if err != nil {
		logger.Error("Unable to sign JWT", err)
		http.Error(w, "Unable to sign JWT", http.StatusInternalServerError)
		return
	}

	response := Response{
		Status:  true,
		Message: "Successfully signed up",
		Data: map[string]any{
			"token": token,
			"id":    userID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
