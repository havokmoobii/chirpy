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
	dbChirps, err := cfg.DB.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}
	
	parsed := uuid.Nil
	s := r.URL.Query().Get("author_id")
	if s != "" {
		parsed, err = uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't parse request author_id", err)
			return
		}
	}

	chirps := make([]Chirp, 0)

	for _, dbChirp := range dbChirps {
		if s != "" && dbChirp.UserID != parsed {
			continue
		}

		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "desc" {
		sortedChirps := make([]Chirp, len(chirps))
		for i := 0; i < len(chirps); i++ {
			sortedChirps[len(chirps) - (i + 1)] = chirps[i] 
		}
		respondWithJSON(w, http.StatusOK, sortedChirps)
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}