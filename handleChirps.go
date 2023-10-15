package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (apiCfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters")
		return
	}

	// Check maximum length
	if len(params.Body) > 140 {
		log.Printf("Message length exceeded: %v", len(params.Body))
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	// Remove bad words
	cleanedMsg := filterBadWords(params.Body)

	// Save it to the store
	chirp, err := apiCfg.DB.CreateChirp(cleanedMsg)
	if err != nil {
		respondWithError(w, 500, "Can't create chirp")
		return
	}

	respondWithJSON(w, 200, databaseChirpToChirp(chirp))
}

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := apiCfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, 500, "Couldn't retrive chirps")
	}

	respondWithJSON(w, 200, databaseChirpsToChirps(dbChirps))
}
