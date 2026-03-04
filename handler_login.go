package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavidMWeaver4/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve password", err)
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
	//get expires in seconds values or set default
	expiresIn := time.Hour
	if params.ExpiresInSeconds != 0 {
		requestedDuration := time.Duration(params.ExpiresInSeconds) * time.Second
		if requestedDuration < expiresIn {
			expiresIn = requestedDuration
		}
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	//final response
	respondWithJSON(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})

}
