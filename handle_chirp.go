package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

	respondWithJSON(w, 201, databaseChirpToChirp(chirp))
}

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := apiCfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, 500, "Couldn't retrive chirps")
	}

	respondWithJSON(w, 200, databaseChirpsToChirps(dbChirps))
}

func (apiCfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := chi.URLParam(r, "chirpID")
	dbChirps, err := apiCfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, 500, "Couldn't retrive chirps")
	}

	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid ID: %v received", chirpIDStr))
		return
	}

	for _, chirp := range dbChirps {
		if chirp.ID == chirpID {
			respondWithJSON(w, 200, databaseChirpToChirp(chirp))
			return
		}
	}

	respondWithError(w, 404, fmt.Sprintf("Chirp with ID: %d not found", chirpID))
}
