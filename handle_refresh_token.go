package main

import (
	"net/http"
	"time"

	"github.com/utphalax/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, bearerErr := auth.GetBearerToken(r.Header)

	if bearerErr != nil {
		responseWithError(w, 401, "Unauthorized", bearerErr)
		return
	}

	// check token in the refresh tokens table
	refreshTokenRecord, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)

	if err != nil {
		responseWithError(w, 401, "Unauthorized", err)
		return
	}

	if refreshTokenRecord.RevokedAt.Valid || refreshTokenRecord.ExpiresAt.Before(time.Now()) {
		responseWithError(w, 401, "Unauthorized", nil)
		return
	}

	// create new JWT token
	jwtToken, err := auth.MakeJWT(refreshTokenRecord.UserID, cfg.jwtSecret, time.Duration(3600)*time.Second)

	if err != nil {
		responseWithError(w, 500, "Failed to create token", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	responseWithJSON(w, 200, response{Token: jwtToken})
}
