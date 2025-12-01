package api

import (
	"net/http"
	"time"
	"strconv"

	"github.com/havokmoobii/chirpy/internal/auth"
)

func (cfg *APIConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Token          string    `json:"token"`
	}
	
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't extract token from header", err)
		return
	}

	tokenExpired, err := cfg.DB.GetRefreshTokenExpired(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user ID from token", err)
		return
	}
	if tokenExpired {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has expired", err)
		return
	}

	tokenRevoked, err := cfg.DB.GetRefreshTokenRevoked(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user ID from token", err)
		return
	}
	if tokenRevoked.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has been revoked", err)
		return
	}

	userID, err := cfg.DB.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user ID from token", err)
		return
	}

	seconds, _ := time.ParseDuration(strconv.Itoa(3600) + "s")

	accessToken, err := auth.MakeJWT(userID, cfg.Secret, seconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Token:        accessToken,
	})
}