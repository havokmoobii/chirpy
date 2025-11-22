package main

import (
	"net/http"
	"time"
	"strconv"


	"fmt"

	"github.com/havokmoobii/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Token          string    `json:"token"`
	}
	
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't extract token from header", err)
		return
	}

	tokenExpired, err := cfg.db.GetRefreshTokenExpired(r.Context(), token)
	if err != nil || tokenExpired {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user ID from token", err)
		return
	}

	tokenRevoked, err := cfg.db.GetRefreshTokenRevoked(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user ID from token", err)
		return
	}

	// Refactor tokenRevoked. Need to get whether or not the value is NULL.

	fmt.Println(tokenRevoked)

	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user ID from token", err)
		return
	}

	seconds, _ := time.ParseDuration(strconv.Itoa(3600) + "s")

	accessToken, err := auth.MakeJWT(userID, cfg.secret, seconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Token:        accessToken,
	})
}