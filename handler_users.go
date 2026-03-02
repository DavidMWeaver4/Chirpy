package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/DavidMWeaver4/Chirpy/internal/database"

)
func (cfg *apiConfig) handlerUsers (w http.ResponseWriter, r *http.Request){

	type parameters struct{
		Email string `json:"email"`
	}
	//decode request body

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	//create DB params
	dbparams := database.CreateUserParams{
		ID: 			uuid.New(),
		CreatedAt: 		time.Now().UTC(),
		UpdatedAt:		time.Now().UTC(),
		Email:			params.Email,
	}

	//create user
	user, err := cfg.db.CreateUser(r.Context(), dbparams)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	//response
	respondWithJSON(w, http.StatusCreated, User{
		ID: 		user.ID,
		CreatedAt:	user.CreatedAt,
		UpdatedAt:	user.UpdatedAt,
		Email:		user.Email,
	})
}
