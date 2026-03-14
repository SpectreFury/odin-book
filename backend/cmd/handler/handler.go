package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/SpectreFury/odin-book/backend/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
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

	var id int32
	err = h.DB.QueryRow(context.Background(), `SELECT id FROM users WHERE email = $1 LIMIT 1`, user.Email).Scan(&id)

	if err != pgx.ErrNoRows {
		http.Error(w, "User already exists, try logging in", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Unable to hash password", err)
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
		return
	}

	result, err := h.DB.Exec(context.Background(), `INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4);`, user.FirstName, user.LastName, user.Email, hashedPassword)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info(result.String())

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "odin-book",
		"sub": "userid",
		"id":  id,
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
			"id":    id,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
