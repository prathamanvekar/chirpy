package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/prathamanvekar/learn-http-servers/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const oneHourInSeconds = 3600
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = oneHourInSeconds
	}

	if params.ExpiresInSeconds > oneHourInSeconds {
		params.ExpiresInSeconds = oneHourInSeconds
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.serverSecret, time.Duration(params.ExpiresInSeconds) * time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failure making the token!", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Token:     token,
		},
	})
}
