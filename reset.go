package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev"{
		respondWithError(w, http.StatusForbidden, "403 Forbidden", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.ResetUsers(r.Context())
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Cannot delete users", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
	return
}
