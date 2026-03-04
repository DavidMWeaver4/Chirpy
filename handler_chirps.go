package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/DavidMWeaver4/Chirpy/internal/auth"
	"github.com/DavidMWeaver4/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140
	type params struct {
		Body string `json:"body"`
	}
	//decode and check if valid
	decoder := json.NewDecoder(r.Body)
	param := params{}
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, 500, "Decoding error", err)
		return
	}
	if len(param.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	//profanity filer
	words := strings.Split(param.Body, " ")
	param.Body = strings.Join(validate_words(words), " ")

	//validations
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}
	validated, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}
	//save to db
	dbParams := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body:      param.Body,
		UserID:    validated,
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, 500, "Failed to save Chirp to DB", err)
		return
	}
	//successful
	respondWithJSON(w, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

// profanity filter helper function
func validate_words(wordstoCheck []string) []string {
	words := wordstoCheck
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" ||
			strings.ToLower(word) == "sharbert" ||
			strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}
	}
	return words
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 404, "Failed to get chirps from database", err)
		return
	}
	var response []Chirp
	for _, chirp := range chirps {
		response = append(response, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, 200, response)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	idS, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, 400, "No chirpID given", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), idS)
	if err != nil {
		respondWithError(w, 404, "Failed to get chirp from database", err)
		return
	}
	respondWithJSON(w, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
