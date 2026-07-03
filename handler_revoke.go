package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/prathamanvekar/learn-http-servers/internal/auth"
	"github.com/prathamanvekar/learn-http-servers/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Token:     refreshToken,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke refresh token", err)
		return
	}

	w.WriteHeader(204)
}
