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
	secr := os.Getenv("SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}
	data := database.New(db)

	const port = "8080"
	const filepathRoot = "."
	apiCfg := apiConfig{
		db:        data,
		platform:  platf,
		jwtSecret: secr,
		PolkaKey:  polkaKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	mux.HandleFunc("GET /api/chirps/", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersAuth)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDelChirps)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handleSetChripyRed)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
