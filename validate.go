package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnError struct {
		Error string `json:"error"`
	}

	type returnValid struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		respBody := returnError{
			Error: "Something went wrong",
		}
		dat, _ := json.Marshal(respBody)
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Message length exceeded: %v", len(params.Body))
		w.WriteHeader(400)
		respBody := returnError{
			Error: "Chirp is too long",
		}
		dat, _ := json.Marshal(respBody)
		w.Write(dat)
		return
	}

	w.WriteHeader(http.StatusOK)
	respBody := returnValid{
		Valid: true,
	}
	dat, _ := json.Marshal(respBody)
	w.Write(dat)
}
