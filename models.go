package main

import "github.com/balazsjuhasz/go-web-server/internal/database"

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func databaseChirpToChirp(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	}
}

func databaseChirpsToChirps(dbChirps []database.Chirp) []Chirp {
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, databaseChirpToChirp(dbChirp))
	}
	return chirps
}
