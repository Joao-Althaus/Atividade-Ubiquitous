package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"users-api-go/internal/user"
)

func main() {
	databaseURL := buildDatabaseURL()
	database, err := connectToDatabase(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	userRepository := user.NewPostgresRepository(database)
	userService := user.NewService(userRepository)
	userController := user.NewController(userService)

	router := http.NewServeMux()
	userController.RegisterRoutes(router)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Go API listening on :8080")
	log.Fatal(server.ListenAndServe())
}

func buildDatabaseURL() string {
	host := env("DB_HOST", "localhost")
	port := env("DB_PORT", "5432")
	name := env("DB_NAME", "users_db")
	username := env("DB_USER", "users")
	password := env("DB_PASSWORD", "users_dev_password")
	sslMode := env("DB_SSLMODE", "disable")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		username,
		password,
		host,
		port,
		name,
		sslMode,
	)
}

func connectToDatabase(databaseURL string) (*pgxpool.Pool, error) {
	context := context.Background()
	database, err := pgxpool.New(context, databaseURL)
	if err != nil {
		return nil, err
	}

	for attempt := 0; attempt < 30; attempt++ {
		if err := database.Ping(context); err == nil {
			return database, nil
		}
		log.Println("waiting for postgres")
		time.Sleep(time.Second)
	}

	database.Close()
	return nil, fmt.Errorf("could not connect to postgres after 30 attempts")
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
