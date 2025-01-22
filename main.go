package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MdSadiqMd/RSS-Scraper/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	connection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	db := database.New(connection)
	apiConfig := apiConfig{
		DB: database.New(connection),
	}

	go startScrapping(db, 10, 5*time.Minute)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", apiConfig.handleCreateUser)
	v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.handleGetUser))
	v1Router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.handleCreateFeed))
	v1Router.Get("/feeds", apiConfig.handleGetFeeds)
	v1Router.Post("/feed_follows", apiConfig.middlewareAuth(apiConfig.handleCreateFeedFollow))
	v1Router.Get("/feed_follows", apiConfig.middlewareAuth(apiConfig.handleGetFeedFollows))
	v1Router.Delete("/feed_follows/{feed_id}", apiConfig.middlewareAuth(apiConfig.handleDeleteFeedFollow))
	v1Router.Get("/posts", apiConfig.middlewareAuth(apiConfig.handleGetPostsForUser))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	log.Printf("Server started on port: %s", PORT)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
