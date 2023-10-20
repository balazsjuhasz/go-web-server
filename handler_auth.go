package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/balazsjuhasz/go-web-server/internal/auth"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiredInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := apiCfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid login credentials")
		return
	}

	token, err := auth.CreateJwtToken("access", strconv.Itoa(user.ID), apiCfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}

	refresh_token, err := auth.CreateJwtToken("refresh", strconv.Itoa(user.ID), apiCfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToAuthenticatedUser(user, token, refresh_token))
}

func (apiCfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	userIdString, err := auth.ValidateJwtToken(token, apiCfg.jwtSecret, "chirpy-refresh")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	tokenRevoked, err := apiCfg.DB.CheckTokenRevoked(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Token revoked checking failed")
		return
	}

	if tokenRevoked {
		respondWithError(w, http.StatusUnauthorized, "Token is revoked")
		return
	}

	token, err = auth.CreateJwtToken("access", userIdString, apiCfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}

	respondWithJSON(w, http.StatusOK, AccessToken{Token: token})
}

func (apiCfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	_, err = auth.ValidateJwtToken(token, apiCfg.jwtSecret, "chirpy-refresh")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	err = apiCfg.DB.RevokeToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't revoke JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
