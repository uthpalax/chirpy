package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/utphalax/chirpy/internal/auth"
	"github.com/utphalax/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirps(w http.ResponseWriter, r *http.Request) {
	

	token, bearerErr := auth.GetBearerToken(r.Header)
	if bearerErr != nil {
		responseWithError(w, 401, "Unauthorized", bearerErr)
		return
	}

	userId, authErr := auth.ValidateJWT(token, cfg.jwtSecret)

	if authErr != nil {
		responseWithError(w, 401, "Unauthorized", authErr)
		return
	}

	type parameters struct {
		Body   string    `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		responseWithError(w, 500, "Something went wrong", err)
		return
	}

	if len(params.Body) > 140 {
		responseWithError(w, 400, "Chirp is too long", err)
		return
	}

	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	getWords := strings.Split(params.Body, " ")

	for _, badWord := range badWords {
		for i, word := range getWords {
			if strings.ToLower(word) == badWord {
				if len(getWords) <= i+1 && getWords[i+1] == "!" {
					continue
				}

				getWords[i] = "****"
			}
		}
	}

	sanitizedBody := strings.Join(getWords, " ")

	now := time.Now()
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Body:      sanitizedBody,
		UserID:    userId,
	})
	if err != nil {
		responseWithError(w, 500, "Cound not create chirp", err)
	}

	type response struct {
		Chirp
	}

	responseWithJSON(w, 201, response{Chirp: Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}})
}
