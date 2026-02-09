package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joaovictornovais/logiscale/pkg/postgres"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")

	ctx := context.Background()
	pool, err := postgres.NewClient(ctx, connStr)

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("LogiScale is running!"))
	})

	port := os.Getenv("PORT")

	log.Println("Server starting on port " + port)
	http.ListenAndServe(":"+port, r)
}
