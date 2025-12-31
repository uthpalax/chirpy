package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/utphalax/chirpy/internal/auth"
	"github.com/utphalax/chirpy/internal/database"
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
		RefreshToken string `json:"refresh_token"`
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

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(expiresInSeconds)*time.Second)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Could not create JWT token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	// store the refreshtoken in the database
	now := time.Now()
	refreshTokenExpiresAt := now.Add(60 * 24 * time.Hour) // 7 days

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		ExpiresAt: refreshTokenExpiresAt,
		RevokedAt: sql.NullTime{Valid: false},
	})

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Could not create refresh token", err)
		return
	}

	responseWithJSON(w, 200, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
		RefreshToken: refreshToken,
	})
}
