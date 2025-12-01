package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/havokmoobii/chirpy/internal/auth"
)

func (cfg *APIConfig) HandlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	parsed, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse request chirp_id", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't extract token from header", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	chirp, err := cfg.DB.GetChirp(r.Context(), parsed)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Cannot delete other user's chirps", err)
		return
	}

	err = cfg.DB.DeleteChirp(r.Context(), parsed)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}