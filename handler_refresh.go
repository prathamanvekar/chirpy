package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/prathamanvekar/learn-http-servers/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token", err)
		return
	}

	refreshTokenRecord, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get refresh token", err)
		return
	}

	isExpired := time.Now().After(refreshTokenRecord.ExpiresAt)
	isRevoked := refreshTokenRecord.RevokedAt.Valid

	if isExpired || isRevoked {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is expired or revoked", errors.New("token expired or revoked"))
		return
	}

	accessToken, err := auth.MakeJWT(
		refreshTokenRecord.UserID,
		cfg.jwtSecret,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
