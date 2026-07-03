package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/prathamanvekar/learn-http-servers/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	// Validate the Polka API Key sent via Authorization header
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
		return
	}
	if apiKey != cfg.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	// Decode webhook event request body
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode", err)
		return
	}

	// Only process 'user.upgraded' events. Respond with 204 for all other event types.
	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	// Upgrade the user in the database. Returns sql.ErrNoRows (404) if user doesn't exist.
	_, err = cfg.db.UpgradeUser(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to upgrade user", err)
		return
	}

	w.WriteHeader(204)
}
