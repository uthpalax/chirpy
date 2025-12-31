package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/utphalax/chirpy/internal/auth"
	"github.com/utphalax/chirpy/internal/database"
)

func (cfg *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, bearerErr := auth.GetBearerToken(r.Header)

	if bearerErr != nil {
		responseWithError(w, 401, "Unauthorized", bearerErr)
		return
	}

	_, err := cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: time.Now(),
		Token:     refreshToken,
	})

	if err != nil {
		responseWithError(w, 500, "Failed to revoke token", err)
		return
	}

	responseWithJSON(w, 204, nil)	
}
