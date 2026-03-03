package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavidMWeaver4/Chirpy/internal/auth"
	"github.com/DavidMWeaver4/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {

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
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Error hashing passsword", err)
	}

	//create DB params
	dbparams := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Email:          params.Email,
		HashedPassword: hashed_password,
	}

	//create user
	user, err := cfg.db.CreateUser(r.Context(), dbparams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	//response
	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
