package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nicwilliams1/rss-aggregator/internals/database"
	"github.com/nicwilliams1/rss-aggregator/internals/utils"
)

func (cfg *apiConfig) handlerFeedsGet(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't fetch feeds")
		return
	}

	var feeds []Feed
	for _, f := range dbFeeds {
		feeds = append(feeds, databaseFeedToFeed(f))
	}

	utils.RespondWithJSON(w, http.StatusOK, feeds)

}

func (cfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	type results struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
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

	feed, err := cfg.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserId:    user.ID,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserId:    user.ID,
		FeedId:    feed.ID,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, results{
		Feed:       databaseFeedToFeed(feed),
		FeedFollow: databaseFeedFollowToFeedFollow(feedFollow),
	})

}
