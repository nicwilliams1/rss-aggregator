package main

import (
	"net/http"

	"github.com/nicwilliams1/rss-aggregator/internals/utils"
)

func (cfg *apiConfig) handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status string `json:"status"`
	}

	utils.RespondWithJSON(w, http.StatusOK, response{Status: "ok"})
}
