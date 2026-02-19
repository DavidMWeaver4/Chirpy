package main
import ("encoding/json"
	"net/http"
)
func (cfg *apiConfig) validate_chirp(w http.ResponseWriter, r *http.Request){
	const maxChirpLength = 140
	type params struct{
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	param := params{}
	err := decoder.Decode(&param)
	if err != nil{
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if len(param.Body) > maxChirpLength{
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	respondWithJSON(w, 200, struct {
    		Valid bool `json:"valid"`
			}{Valid: true})

}
