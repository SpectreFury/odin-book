package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SpectreFury/odin-book/backend/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Response struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Body    map[string]any `json:"body"`
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

	rows, err := h.DB.Query(context.Background(), `INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4);`, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println(rows)

	response := Response{
		Status:  true,
		Message: "Successfully signed up",
		Body:    map[string]any{},
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
