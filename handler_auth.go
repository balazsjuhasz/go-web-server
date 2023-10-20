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

	token, err := auth.CreateJwtToken(auth.TokenTypeAccess, strconv.Itoa(user.ID), apiCfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}

	refresh_token, err := auth.CreateJwtToken(auth.TokenTypeRefresh, strconv.Itoa(user.ID), apiCfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToAuthenticatedUser(user, token, refresh_token))
}

func (apiCfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	userIdString, err := auth.ValidateJwtToken(refreshToken, apiCfg.jwtSecret, auth.TokenTypeRefresh)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	isTokenRevoked, err := apiCfg.DB.IsTokenRevoked(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Token revoked checking failed")
		return
	}

	if isTokenRevoked {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is revoked")
		return
	}

	accessToken, err := auth.CreateJwtToken(auth.TokenTypeAccess, userIdString, apiCfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}

	respondWithJSON(w, http.StatusOK, AccessToken{Token: accessToken})
}

func (apiCfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	_, err = auth.ValidateJwtToken(refreshToken, apiCfg.jwtSecret, auth.TokenTypeRefresh)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	err = apiCfg.DB.RevokeToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't revoke JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
