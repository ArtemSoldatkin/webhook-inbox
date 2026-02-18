package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	routev1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/routes/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/config"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	deliveryengine "github.com/ArtemSoldatkin/webhook-inbox/internal/delivery_engine"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/sirupsen/logrus"
)

// createDatabasePool creates a new database connection pool.
func createDatabasePool(user string, password string, url string, port int, db string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, url, port, db))
	if err != nil {
		logrus.WithError(err).Fatal("Unable to connect to database")
	}
	return pool
}

// LogrusLogger is a middleware that logs incoming HTTP requests using logrus, including the method, path, and duration of the request.
func LogrusLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        logrus.WithFields(logrus.Fields{
            "method": r.Method,
            "path":   r.URL.Path,
            "duration": time.Since(start),
        }).Info("Handled request")
    })
}


// main is the entry point of the application.
func main() {
	config := config.LoadConfig()
	r := chi.NewRouter()
	r.Use(LogrusLogger)
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

	go deliveryengine.Start(service, 300 * time.Millisecond) // TODO make this configurable
	go deliveryengine.StartRecoveryEngine(service, 5 * time.Minute) // TODO make this configurable

	r.Mount("/api/v1", routev1.V1Router(service))

	err := http.ListenAndServe(fmt.Sprintf(":%s", config.ApiPort), r)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}
}
