package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nicwilliams1/rss-aggregator/internals/database"
	"github.com/nicwilliams1/rss-aggregator/internals/utils"
)

func (cfg *apiConfig) handlerFeedFollowsCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	ctx := r.Context()
	user := ctx.Value("user").(database.User)

	if len(user.ID) == 0 {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't extract user from context")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	feed, err := cfg.DB.GetFeed(r.Context(), params.FeedId)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't check if feed exists")
		return
	}

	if len(feed.ID) == 0 {
		utils.RespondWithError(w, http.StatusInternalServerError, "Feed does not exist")
		return
	}

	dbFeedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedId:    feed.ID,
		UserId:    user.ID,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(dbFeedFollow))

}

func (cfg *apiConfig) handlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	if _, ok := ctx.Value("user").(database.User); !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't extract user from context")
		return
	}

	idStr := chi.URLParam(r, "feedFollowID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid feed follow uuid")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't delete feed follow")
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (cfg *apiConfig) handlerFeedFollowsGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(database.User)

	if len(user.ID) == 0 {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't get feed follows for user")
		return
	}

	dbFeedFollows, err := cfg.DB.GetFeedFollowsByUserId(r.Context(), user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't extract user from context")
		return
	}

	var feedsFollows []FeedFollow
	for _, f := range dbFeedFollows {
		feedsFollows = append(feedsFollows, databaseFeedFollowToFeedFollow(f))
	}

	utils.RespondWithJSON(w, http.StatusOK, feedsFollows)

}
