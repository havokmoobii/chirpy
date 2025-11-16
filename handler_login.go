package main

import(
	"net/http"
	"encoding/json"
	"strings"
	"time"

	"fmt"

	"github.com/google/uuid"
	"github.com/havokmoobii/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string  `json:"email"`
		Password         string  `json:"password"`
		ExpiresInSeconds int     `json:"expires_in_seconds"`
	}
	type returnVals struct {
		ID             uuid.UUID `json:"id"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Email          string    `json:"email"`
		Token          string    `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 3600 {
		params.ExpiresInSeconds = 3600
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	seconds, _ := time.ParseDuration(string(params.ExpiresInSeconds) + "s")
	fmt.Println(string(params.ExpiresInSeconds) + "s")

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't compare hashed password", err)
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, seconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}