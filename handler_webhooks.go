package main

import (
	"encoding/json"
	"net/http"

	"github.com/DavidMWeaver4/Chirpy/internal/auth"
	"github.com/DavidMWeaver4/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleSetChripyRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	//decode
	var param parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, 500, "Decoding error", err)
		return
	}
	//check if the correct condition is being called
	if param.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	//get apikey
	ApiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	if ApiKey != cfg.PolkaKey {
		w.WriteHeader(401)
		return
	}
	user, err := cfg.db.GetUser(r.Context(), param.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if param.Data.UserID != user.ID {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = cfg.db.SetUsersChirpyRed(r.Context(), database.SetUsersChirpyRedParams{
		IsChirpyRed: true,
		ID:          param.Data.UserID,
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
