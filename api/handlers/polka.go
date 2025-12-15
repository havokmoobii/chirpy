package api

import (
	"net/http"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/havokmoobii/chirpy/internal/auth"
)

func (cfg *APIConfig) HandlerUsersUpgradeToRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't extract API key from header", err)
		return
	}

	if apiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API Key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not parse user id", err)
		return
	}

	err = cfg.DB.UpgradeUser(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Could not upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}