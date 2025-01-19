package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

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
	apiConfig := apiConfig{
		DB: database.New(connection),
	}

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
