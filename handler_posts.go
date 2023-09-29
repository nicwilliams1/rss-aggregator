package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/nicwilliams1/rss-aggregator/internals/database"
	"github.com/nicwilliams1/rss-aggregator/internals/utils"
)

func (cfg *apiConfig) handlerPostsGetByUser(w http.ResponseWriter, r *http.Request) {

	defaultLimit := "10"

	ctx := r.Context()
	user := ctx.Value("user").(database.User)

	if len(user.ID) == 0 {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't extract user from context")
		return
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = defaultLimit
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't parse limit query parameter")
		return
	}

	feedFollows, err := cfg.DB.GetFeedFollowsByUserId(r.Context(), user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't fetch feed follows for user")
		return
	}

	feed_ids := make([]uuid.UUID, 0)
	for _, ffs := range feedFollows {
		feed_ids = append(feed_ids, ffs.FeedId)
	}

	dbPosts, err := cfg.DB.GetPostsByUser(r.Context(), feed_ids, limitInt)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't fetch posts for user")
		return
	}

	posts := []Post{}
	for _, p := range dbPosts {
		posts = append(posts, databasePostToPost(p))
	}
	utils.RespondWithJSON(w, http.StatusOK, posts)

}
