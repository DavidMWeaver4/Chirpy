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

func (cfg *apiConfig) handlerUsersAuth(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	//decode
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	//get access token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}
	//get the User's ID
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user", err)
		return
	}
	//hash the new password
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 401, "Error hashing passsword", err)
		return
	}
	//store to DB
	err = cfg.db.ChangeUsersEmailAndPassword(r.Context(), database.ChangeUsersEmailAndPasswordParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, 401, "Error changing password", err)
		return
	}
	//grab user
	user, err := cfg.db.GetUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, 401, "Error fetching user", err)
		return
	}
	//sucess response
	respondWithJSON(w, 200, User{
		ID:        userID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now().UTC(),
		Email:     params.Email,
	})

}
