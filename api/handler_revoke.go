package api

import (
	"net/http"

	"github.com/havokmoobii/chirpy/internal/auth"
)

func (cfg *APIConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't extract token from header", err)
		return
	}

	err = cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}