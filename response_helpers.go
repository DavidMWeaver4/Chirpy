package main
import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string){
	respondWithJSON(w, code, map[string]string{"error": msg})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	dat, err := json.Marshal(payload)
	if err != nil{
	 	w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
