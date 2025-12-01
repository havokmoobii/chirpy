package api

import(
	"net/http"
	"encoding/json"
	"strings"
	"time"
	"strconv"

	"github.com/google/uuid"
	"github.com/havokmoobii/chirpy/internal/auth"
	"github.com/havokmoobii/chirpy/internal/database"
)

func (cfg *APIConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string  `json:"email"`
		Password         string  `json:"password"`
	}
	type returnVals struct {
		ID             uuid.UUID `json:"id"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Email          string    `json:"email"`
		Token          string    `json:"token"`
		RefreshToken    string    `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), params.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't compare hashed password", err)
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	seconds, _ := time.ParseDuration(strconv.Itoa(3600) + "s")

	accessToken, err := auth.MakeJWT(user.ID, cfg.Secret, seconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make JWT", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}