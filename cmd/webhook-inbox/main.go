package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	routev1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/routes/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/config"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
)

// createDatabasePool creates a new database connection pool.
func createDatabasePool(user string, password string, url string, port int, db string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, url, port, db))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	return pool
}

// main is the entry point of the application.
func main() {
	config := config.LoadConfig()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	dbPool := createDatabasePool(
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	defer dbPool.Close()
	queries := db.New(dbPool)
	service := service.NewService(queries)


	r.Mount("/api/v1", routev1.V1Router(service))

	err := http.ListenAndServe(fmt.Sprintf(":%s", config.ApiPort), r)
	if err != nil {
		log.Fatal(err)
	}
}
