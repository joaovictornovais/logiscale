package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/joaovictornovais/logiscale/internal/handler"
	repository "github.com/joaovictornovais/logiscale/internal/repository/postgres"
	"github.com/joaovictornovais/logiscale/internal/service"
	pgPkg "github.com/joaovictornovais/logiscale/pkg/postgres"
	"github.com/joaovictornovais/logiscale/pkg/redis"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()
	pool, err := pgPkg.NewClient(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	rdb, err := redis.NewClient(redisHost + ":" + redisPort)
	if err != nil {
		log.Fatalf("error while connecting to redis: %v", err)
	}
	defer rdb.Close()
	defer pool.Close()

	driverRepo := repository.NewDriverRepository(pool)
	driverService := service.NewDriverService(driverRepo)
	driverHandler := handler.NewDriverHandler(driverService)

	locationRepo := repository.NewLocationRepository(pool)
	ingestionService := service.NewIngestionService(locationRepo, rdb)
	ingestionHandler := handler.NewIngestionHandler(ingestionService)

	dispatchService := service.NewDispatchService(rdb)
	dispatchHandler := handler.NewDispatchHandler(dispatchService)

	defer ingestionService.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("LogiScale is running!"))
	})

	r.Post("/drivers", driverHandler.CreateDriver)
	r.Get("/drivers/{id}", driverHandler.GetByID)
	r.Post("/drivers/{id}/locations", ingestionHandler.HandleIngest)

	r.Get("/drivers/nearest", dispatchHandler.FindNearest)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on port " + port)
	http.ListenAndServe(":"+port, r)
}
