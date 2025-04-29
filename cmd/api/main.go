package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"targeting-engine/configs"
	"targeting-engine/internal/handlers"
	"targeting-engine/internal/middleware"
	"targeting-engine/internal/repository"
	"targeting-engine/internal/service"
)

func main() {
	// Get our settings
	settings := configs.NewConfig()
	settings.LoadFromEnv()

	// Set up PostgreSQL connection
	ctx := context.Background()

	postgresStore, err := repository.NewPostgresRepository(ctx, settings.Database.PostgresURI)
	if err != nil {
		log.Fatalf("Dang! Can't connect to PostgreSQL: %v", err)
	}
	defer func() {
		log.Println("Closing PostgreSQL connection...")
		if err := postgresStore.Close(ctx); err != nil {
			log.Printf("Whoops! Error closing PostgreSQL: %v", err)
		}
	}()

	log.Println("Adding some test data to PostgreSQL...")
	if err := postgresStore.InitTestData(ctx); err != nil {
		log.Printf("Hmm, couldn't add test data: %v", err)
	}

	campaignMatcher := service.NewTargetingService(postgresStore)
	campaignHandler := handlers.NewDeliveryHandler(campaignMatcher)

	router := http.NewServeMux()

	router.Handle("/v1/delivery", campaignHandler)
	if settings.EnableHealthCheck {
		router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	}
	handler := middleware.LoggingMiddleware(router)

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(settings.Port),
		Handler: handler,
	}

	go func() {
		log.Printf("Starting server on port %d", settings.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Uh oh! Server couldn't start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server had to force quit: %v", err)
	}

	log.Println("Server shut down nicely")
}
