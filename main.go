package main

import (
	"log"
	"net/http"
	"os"

	"github.com/balazsjuhasz/go-web-server/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	filepathRoot := "."

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	r := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/metrics", apiCfg.handlerMetrics)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Post("/chirps", apiCfg.handlerCreateChirp)
	apiRouter.Get("/chirps", apiCfg.handlerGetChirps)
	r.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerAdminMetrics)
	r.Mount("/admin", adminRouter)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + portString,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, portString)
	log.Fatal(srv.ListenAndServe())
}
