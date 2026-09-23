package main

import (
	"fmt"
	"net/http"

	"github.com/alexandre-j95/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't parse chirpID: %v", err))
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't find token: %s", err))
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't validate token: %s", err))
		return
	}

	chirpObj, err := cfg.DB.GetChirpByID(r.Context(), chirpID)
		if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't retrieve chirp: %v", err))
			return
		}

	if chirpObj.UserID != userID {
		respondWithError(w, http.StatusForbidden, "User can't delete this chirp")
		return
	}

	err = cfg.DB.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't delete chirp: %v", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
