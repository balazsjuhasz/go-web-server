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
	jwtSecret      string
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not found in the environment")
	}

	filepathRoot := "."

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
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
	apiRouter.Get("/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	apiRouter.Delete("/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	apiRouter.Post("/users", apiCfg.handlerCreateUser)
	apiRouter.Put("/users", apiCfg.handlerUpdateUser)

	apiRouter.Post("/login", apiCfg.handlerLogin)

	apiRouter.Post("/refresh", apiCfg.handlerRefreshToken)
	apiRouter.Post("/revoke", apiCfg.handlerRevokeToken)

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
