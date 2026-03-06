package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavidMWeaver4/Chirpy/internal/auth"
	"github.com/DavidMWeaver4/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	//decode request body

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	//get user
	user, err := cfg.db.GetUserFromEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve user", err)
		return
	}
	//check password
	isValid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Error checking password", err)
		return
	}
	//check if password is valid
	if !isValid {
		respondWithError(w, 401, "Incorrect email or password", nil)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	refreshedTokenString := auth.MakeRefreshToken()
	refresh_token, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshedTokenString,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	//final response
	respondWithJSON(w, 200, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh_token.Token,
		IsChirpyRed:  user.IsChirpyRed,
	})

}
