package main

import (
	"os"
	"log"
	"net/http"
	"sync/atomic"
	"database/sql"

	"github.com/havokmoobii/chirpy/api/handlers"
	"github.com/havokmoobii/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const filepathRoot = "./public"
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening database: %s", err)
		return
	}
	dbQueries := database.New(db)

	apiCfg := api.APIConfig{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
		Platform:       os.Getenv("PLATFORM"),
		Secret:         os.Getenv("SECRET"),
	}

	mux := http.NewServeMux()

	fsHandler := apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", (http.FileServer(http.Dir(filepathRoot)))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", api.HandlerReadiness)
	mux.HandleFunc("POST /api/users", apiCfg.HandlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.HandlerUsersUpdate)
	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandlerRevoke)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandlerChirpsGet)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.HandlerChirpsDelete)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.HandlerUsersUpgradeToRed)

	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)
	
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}