package main
import ("encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/DavidMWeaver4/Chirpy/internal/database"
)
func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request){
	const maxChirpLength = 140
	type params struct{
		Body string `json:"body"`
		UserID string`json:"user_id"`
	}
	//decode and check if valid
	decoder := json.NewDecoder(r.Body)
	param := params{}
	err := decoder.Decode(&param)
	if err != nil{
		respondWithError(w, 500, "Decoding error", err)
		return
	}
	if len(param.Body) > maxChirpLength{
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	//profanity filer
	words := strings.Split(param.Body, " ")
	param.Body = strings.Join(validate_words(words), " ")

	//create and save to database
	params_user_id, err := uuid.Parse(param.UserID)
	if err != nil{
		respondWithError(w, 400, "Invalid user_id", err)
		return
	}
	dbParams := database.CreateChirpParams{
		ID:				uuid.New(),
		CreatedAt:		time.Now().UTC(),
		UpdatedAt:		time.Now().UTC(),
		Body:			param.Body,
		UserID:			params_user_id,
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), dbParams)
	if err != nil{
		respondWithError(w, 500, "Failed to save Chirp to DB", err )
		return
	}
	//successful
	respondWithJSON(w, 201, Chirp{
		ID:		chirp.ID,
		CreatedAt:	chirp.CreatedAt,
		UpdatedAt:	chirp.UpdatedAt,
		Body:		chirp.Body,
		UserId:		chirp.UserID,
	})
	return

}

//profanity filter helper function
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
