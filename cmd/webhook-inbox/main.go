package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/config"
)


func main() {
	config := config.LoadConfig()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	err := http.ListenAndServe(fmt.Sprintf(":%s", config.ApiPort), r)
	if err != nil {
		log.Fatal(err)
	}
}
