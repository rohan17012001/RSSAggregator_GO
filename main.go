package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rohan17012001/rssaggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// fmt.Println("Hello, World!")
	// feed, err := urlToFeed("https://www.wagslane.dev/index.xml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(feed)
	godotenv.Load()
	portNo := os.Getenv("PORT")
	if portNo == "" {
		log.Fatal("PORT environment variable not set")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found in environment variables")
	}
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}
	queries := database.New(conn)
	apiCfg := &apiConfig{
		DB: queries,
	}
	go startScraping(queries, 10, time.Minute)
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/user", apiCfg.handlerCreateUser)
	v1Router.Get("/user", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feed", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feed", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollows))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))
	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portNo,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Port number is", portNo)
}
