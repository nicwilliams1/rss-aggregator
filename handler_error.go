package main

import (
	"net/http"

	"github.com/nicwilliams1/rss-aggregator/internals/utils"
)

func (cfg *apiConfig) handlerError(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
