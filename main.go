package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/DavidMWeaver4/Chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	//database and env load
	godotenv.Load()
	platf := os.Getenv("PLATFORM")
	db_url := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}
	data := database.New(db)
	if err != nil {
		log.Fatal(err)
	}

	const port = "8080"
	const filepathRoot = "."
	apiCfg := apiConfig{
		db:       data,
		platform: platf,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	//mux.HandleFunc("POST /api/validate_chirp", apiCfg.validate_chirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
