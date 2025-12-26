package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/utphalax/chirpy/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int   `json:"expires_in_seconds"` // in seconds
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Email password combination did not match", err)
		return
	}

	authenticated, err := auth.ComparePasswordAndHash(params.Password, user.HashedPassword)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Something went wrong", err)
		return
	}

	if !authenticated {
		responseWithError(w, 401, "Email password combination did not match", err)
		return
	}

	expiresInSeconds := 3600 // default 1 hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < expiresInSeconds {
		expiresInSeconds = params.ExpiresInSeconds
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Invalid user ID", err)
		return
	}

	token, err := auth.MakeJWT(userID, cfg.jwtSecret, time.Duration(expiresInSeconds)*time.Second)

	responseWithJSON(w, 200, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}
