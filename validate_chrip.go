package main
import ("encoding/json"
	"net/http"
	"strings"
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
	words := strings.Split(param.Body, " ")
	param.Body = strings.Join(validate_words(words), " ")
	respondWithJSON(w, 200, struct {
    		CleanedBody string `json:"cleaned_body"`
		}{CleanedBody: param.Body})


}
func validate_words(wordstoCheck []string) []string{
	words := wordstoCheck
	for i, word := range words{
		if (strings.ToLower(word) == "kerfuffle" ||
		strings.ToLower(word) == "sharbert" ||
		strings.ToLower(word) == "fornax"){
			words[i] = "****"
		}
	}
	return words
}
