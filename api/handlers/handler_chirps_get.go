package api

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *APIConfig) HandlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	parsed, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse request chirp_id", err)
		return
	}

	chirp, err := cfg.DB.GetChirp(r.Context(), parsed)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *APIConfig) HandlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	if s != "" {
		parsed, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't parse request author_id", err)
			return
		}

		dbChirps, err := cfg.DB.GetChirpsFromID(r.Context(), parsed)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}

		chirps := make([]Chirp, len(dbChirps))

		for i, dbChirp := range dbChirps {
			chirps[i].ID = dbChirp.ID
			chirps[i].CreatedAt = dbChirp.CreatedAt
			chirps[i].UpdatedAt = dbChirp.UpdatedAt
			chirps[i].Body = dbChirp.Body
			chirps[i].UserID = dbChirp.UserID
		}

		respondWithJSON(w, http.StatusOK, chirps)
		return
	}
	
	dbChirps, err := cfg.DB.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	chirps := make([]Chirp, len(dbChirps))

	for i, dbChirp := range dbChirps {
		chirps[i].ID = dbChirp.ID
		chirps[i].CreatedAt = dbChirp.CreatedAt
		chirps[i].UpdatedAt = dbChirp.UpdatedAt
		chirps[i].Body = dbChirp.Body
		chirps[i].UserID = dbChirp.UserID
	}

	respondWithJSON(w, http.StatusOK, chirps)
}