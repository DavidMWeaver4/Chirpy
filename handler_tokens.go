package main

import (
	"net/http"
	"time"

	"github.com/DavidMWeaver4/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "Error getting refresh token", err)
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}
	newJWTAccessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}
	type response struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, 200, response{
		Token: newJWTAccessToken,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "Error getting token", err)
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 400, "Error revoking token", err)
		return
	}
	w.WriteHeader(204)

}
