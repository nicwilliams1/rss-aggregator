package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nicwilliams1/rss-aggregator/internals/database"
	"github.com/nicwilliams1/rss-aggregator/internals/utils"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create user")
	}

	utils.RespondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (cfg *apiConfig) handlerUsersGetByApiKey(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := ctx.Value("user").(database.User)

	if len(user.ID) == 0 {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't extract user from context")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
