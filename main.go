package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nicwilliams1/rss-aggregator/internals/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load()
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	conn := os.Getenv("CONN")

	// connect to db
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	r.Use(cors.Handler)

	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", apiCfg.handlerReadiness)
	v1Router.Get("/err", apiCfg.handlerError)
	v1Router.Post("/users", apiCfg.handlerUsersCreate)
	v1Router.Get("/feeds", apiCfg.handlerFeedsGet)

	// apikey endpoints (require ApiKey in Authorization header)
	v1Router.Group(func(v1Router chi.Router) {
		v1Router.Use(apiCfg.middlewareAuth)
		v1Router.Get("/users", apiCfg.handlerUsersGetByApiKey)
		v1Router.Post("/feeds", apiCfg.handlerFeedsCreate)
		v1Router.Post("/feed_follows", apiCfg.handlerFeedFollowsCreate)
		v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.handlerFeedFollowsDelete)
		v1Router.Get("/feed_follows", apiCfg.handlerFeedFollowsGet)
		v1Router.Get("/posts", apiCfg.handlerPostsGetByUser)
	})

	r.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// const collectionConcurrency = 10
	// const collectionInterval = time.Minute
	// go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())

}
