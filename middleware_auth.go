package main

import (
	"net/http"

	"github.com/MdSadiqMd/RSS-Scraper/internal/auth"
	"github.com/MdSadiqMd/RSS-Scraper/internal/database"
)

type authHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, "Unauthorized")
			return
		}
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 403, "Unauthorized")
			return
		}
		handler(w, r, user)
	}
}
