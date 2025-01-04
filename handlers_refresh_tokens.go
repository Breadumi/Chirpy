package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Breadumi/Chirpy/internal/auth"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, req *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	refreshTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error retrieving refresh token")
		respondWithError(w, http.StatusBadRequest, "error retrieving refresh token")
		return
	}

	dbRefreshToken, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshTokenString)
	if err != nil {
		log.Printf("refresh token does not exist: %s", err)
		respondWithError(w, http.StatusUnauthorized, "refresh token does not exist")
		return
	} else if time.Now().After(dbRefreshToken.ExpiresAt) || dbRefreshToken.RevokedAt.Valid {
		log.Printf("refresh token expired: %s", err)
		respondWithError(w, http.StatusUnauthorized, "refresh token expired")
		return
	}

	jwtTokenString, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("error creating JWT: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error creating JWT")
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: jwtTokenString,
	})

}

func (cfg *apiConfig) revoke(w http.ResponseWriter, req *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error retrieving refresh token")
		respondWithError(w, http.StatusBadRequest, "error retrieving refresh token")
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), refreshTokenString)
	if err != nil {
		log.Printf("error revoking token in database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error revoking token in database")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
