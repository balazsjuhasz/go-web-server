package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/balazsjuhasz/go-web-server/internal/auth"
	"github.com/balazsjuhasz/go-web-server/internal/database"
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	userIdString, err := auth.ValidateJwtToken(token, apiCfg.jwtSecret, "chirpy-access")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user id")
		return
	}

	// Check maximum length
	if len(params.Body) > 140 {
		log.Printf("Message length exceeded: %v", len(params.Body))
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Remove bad words
	cleanedMsg := filterBadWords(params.Body)

	// Save it to the store
	chirp, err := apiCfg.DB.CreateChirp(cleanedMsg, userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseChirpToChirp(chirp))
}

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")
	sorting := r.URL.Query().Get("sort")

	dbChirps, err := apiCfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrive chirps")
		return
	}

	if sorting == "desc" {
		sort.Slice(dbChirps, func(i, j int) bool {
			return dbChirps[i].ID > dbChirps[j].ID
		})
	} else {
		sort.Slice(dbChirps, func(i, j int) bool {
			return dbChirps[i].ID < dbChirps[j].ID
		})
	}

	if authorIDStr == "" {
		respondWithJSON(w, http.StatusOK, databaseChirpsToChirps(dbChirps))
		return
	}

	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse author_id")
		return
	}
	filteredDbChirps := []database.Chirp{}
	for _, dbChirp := range dbChirps {
		if dbChirp.AuthorID == authorID {
			filteredDbChirps = append(filteredDbChirps, dbChirp)
		}
	}

	respondWithJSON(w, http.StatusOK, databaseChirpsToChirps(filteredDbChirps))
}

func (apiCfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := chi.URLParam(r, "chirpID")
	dbChirps, err := apiCfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrive chirps")
	}

	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid ID: %v received", chirpIDStr))
		return
	}

	for _, chirp := range dbChirps {
		if chirp.ID == chirpID {
			respondWithJSON(w, http.StatusOK, databaseChirpToChirp(chirp))
			return
		}
	}

	respondWithError(w, http.StatusNotFound, fmt.Sprintf("Chirp with ID: %d not found", chirpID))
}

func (apiCfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid ID: %v received", chirpIDStr))
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	userIdString, err := auth.ValidateJwtToken(token, apiCfg.jwtSecret, "chirpy-access")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user id")
		return
	}

	dbChirp, err := apiCfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	if dbChirp.AuthorID != userId {
		respondWithError(w, http.StatusForbidden, "Not allowed")
		return
	}

	err = apiCfg.DB.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
