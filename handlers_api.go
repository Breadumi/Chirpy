package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Breadumi/Chirpy/internal/auth"
	"github.com/Breadumi/Chirpy/internal/database"
	"github.com/google/uuid"
)

func readinessEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type response struct {
		Chirp
	}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("unable to retrieve token: %s", err)
		respondWithError(w, http.StatusBadRequest, "unable to retrieve token from request")
		return
	} else if len(tokenString) == 64 {
		log.Printf("invalid token - JWT required: %s", err)
		respondWithError(w, http.StatusUnauthorized, "invalid token - JWT required")
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		log.Printf("invalid token:%s\n %s", tokenString, err)
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	params := parameters{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "error decoding parameters")
		return
	}
	//prettyprint(params) // debugging line

	// check if length is acceptable
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	})
	//prettyprint(chirp) // debugging line
	if err != nil {
		log.Printf("Error querying database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error querying database")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      cleanText(chirp.Body),
			UserID:    chirp.UserID,
		},
	})

}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, req *http.Request) {

	chirpsSQLC, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		log.Printf("error querying database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error querying database")
		return
	}

	var chirps []Chirp

	for _, c := range chirpsSQLC {
		chirps = append(chirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, req *http.Request) {

	type response struct {
		Chirp
	}

	chirpID := req.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		log.Printf("error converting JSON to UUID: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error converting JSON to UUID")
		return
	}

	dbChirp, err := cfg.db.GetChirp(req.Context(), chirpUUID)
	if err != nil {
		log.Printf("error retrieving chirp: %s", err)
		respondWithError(w, http.StatusNotFound, "invalid id")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		},
	})

}

func (cfg *apiConfig) createUser(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters")
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("error hashing password: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "error hashing password")
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})
	//prettyprint(user) // debugging line
	if err != nil {
		log.Printf("error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error retrieving user")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})

}

func (cfg *apiConfig) login(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("error decoding paramaters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters")
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		log.Printf("error retrieving user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error retrieving user")
		return
	}

	if auth.CheckPassword(params.Password, dbUser.HashedPassword) != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	tokenString, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("error in MakeJWT: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error creating JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("error creating refresh token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token")
	}

	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		log.Printf("error storing refresh token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error storing refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:           dbUser.ID,
			CreatedAt:    dbUser.CreatedAt,
			UpdatedAt:    dbUser.UpdatedAt,
			Email:        dbUser.Email,
			Token:        tokenString,
			RefreshToken: refreshToken,
			IsChirpyRed:  dbUser.IsChirpyRed,
		},
	})

}

func (cfg *apiConfig) updateLogin(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	jwtString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error retrieving refresh token")
		respondWithError(w, http.StatusUnauthorized, "error retrieving refresh token")
		return
	}

	userID, err := auth.ValidateJWT(jwtString, cfg.secret)
	if err != nil {
		log.Printf("invalid jwt: %s", err)
		respondWithError(w, http.StatusUnauthorized, "malformed or expired jwt")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("malformed request body")
		respondWithError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	hashPass, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Unable to hash password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "unable to hash password")
		return
	}

	dbUser, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashPass,
	})
	if err != nil {
		log.Printf("user does not exist")
		respondWithError(w, http.StatusBadRequest, "user does not exist")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
	})

}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, req *http.Request) {

	jwtString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error retrieving token")
		respondWithError(w, http.StatusUnauthorized, "error retrieving token")
		return
	}

	userUUID, err := auth.ValidateJWT(jwtString, cfg.secret)
	if err != nil {
		log.Printf("invalid jwt: %s", err)
		respondWithError(w, http.StatusUnauthorized, "malformed or expired jwt")
		return
	}

	chirpID := req.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		log.Printf("error parsing chirpID: %s", err)
		respondWithError(w, http.StatusBadRequest, "malformed chirp ID")
		return
	}

	_, err = cfg.db.DeleteChirp(req.Context(), database.DeleteChirpParams{
		ID:     chirpUUID,
		UserID: userUUID,
	})

	if err != nil {
		log.Printf("error deleting chirp or not found: %s", err)
		respondWithError(w, http.StatusForbidden, "chirp not found")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")

}
