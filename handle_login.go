package main

import (
	"encoding/json"
	"net/http"

	"github.com/utphalax/chirpy/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
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

	responseWithJSON(w, 200, response{User: User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}})
}
