package main

import (
	"net/http"
	//"time"
	//"strconv"

	"github.com/havokmoobii/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	_, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't extract token from header", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}