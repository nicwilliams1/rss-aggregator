package main

import (
	"context"
	"net/http"

	"github.com/nicwilliams1/rss-aggregator/internals/auth"
	"github.com/nicwilliams1/rss-aggregator/internals/utils"
)

func (cfg *apiConfig) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		api_key, err := auth.GetBearerToken(r.Header, "ApiKey")
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Missing ApiKey Authorization header")
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), api_key)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't find user")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
