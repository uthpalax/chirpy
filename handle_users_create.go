package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/utphalax/chirpy/internal/auth"
	"github.com/utphalax/chirpy/internal/database"
)

type User struct {
	ID        string `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Coud not hash password", err)
		return
	}

	now := time.Now()
	payload := database.CreateUserParams{
		ID:             uuid.New().String(),
		CreatedAt:      now,
		UpdatedAt:      now,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.CreateUser(r.Context(), payload)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	responseWithJSON(w, http.StatusCreated, response{User: User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}})
}
