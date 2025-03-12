package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/divakaivan/rssagg/internal/database"
	"github.com/divakaivan/rssagg/internal/handlers"
	"github.com/divakaivan/rssagg/internal/rss"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("[Error] PORT not found in env")
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("[Error] dbUrl not found in env")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("[Error] Cannot connect to database: ", err)
	}

	db := database.New(conn)
	queries := database.New(conn)
	api := handlers.NewAPI(queries)
	// need to start it on its own goroutine
	// so it doesn't block the main thread
	go rss.StartScraping(db, 5, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(api.LogToDbMiddleware)

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", api.HandlerCreateUser)
	v1Router.Get("/err", handlers.HandlerErr)
	v1Router.Post("/users", api.HandlerCreateUser)
	v1Router.Get("/users", api.AuthMiddleware(api.HandlerGetUser))

	v1Router.Post("/feeds", api.AuthMiddleware(api.HandlerCreateFeed))
	v1Router.Get("/feeds", api.HandlerGetFeeds)

	v1Router.Get("/posts", api.AuthMiddleware(api.HandlerGetPostsForUser))

	v1Router.Post("/feed_follows", api.AuthMiddleware(api.HandlerCreateFeedFollow))
	v1Router.Get("/feed_follows", api.AuthMiddleware(api.HandlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", api.AuthMiddleware(api.HandlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("[Info] Server listening on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("[Error] Server failed to start: ", err)
	}
}
