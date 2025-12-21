package main

import (
	"net/http"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		responseWithError(w, 500, "Could not fetch chirps", err)
		return
	}

	chirps := make([]Chirp, 0, len(dbChirps))

	for _, u := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			Body:      u.Body,
			UserID:    u.UserID,
		})
	}

	responseWithJSON(w, 200, chirps)
}
