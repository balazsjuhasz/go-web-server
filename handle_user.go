package main

import (
	"encoding/json"
	"net/http"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters")
		return
	}

	// Save it to the store
	user, err := apiCfg.DB.CreateUser(params.Email)
	if err != nil {
		respondWithError(w, 500, "Can't create user")
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))
}