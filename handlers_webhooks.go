package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Breadumi/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) registerChirpyRed(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil || apiKey != cfg.polkaKey {
		log.Printf("%v", err)
		respondWithError(w, http.StatusUnauthorized, "invalid API key")
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("error decoding JSON: %s", err)
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, "")
		return
	}

	_, err = cfg.db.UpgradeUser(req.Context(), params.Data.UserID)
	if err != nil {
		log.Printf("error finding user: %s", err)
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")

}
