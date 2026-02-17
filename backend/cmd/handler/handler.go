package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AuthHandler struct{}

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

	fmt.Println(user)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}
