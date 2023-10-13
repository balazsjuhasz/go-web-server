package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	filepathRoot := "."

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	r := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)
	r.Get("/healthz", handlerReadiness)
	r.Get("/metrics", apiCfg.handlerMetrics)
	r.Get("/reset", apiCfg.handlerReset)
	r.Post("/api/validate_chirp", apiCfg.handlerValidateChirp)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + portString,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, portString)
	log.Fatal(srv.ListenAndServe())
}
