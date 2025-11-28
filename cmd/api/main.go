package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielyang21/GoBattleServer/internal/database"
	"github.com/danielyang21/GoBattleServer/internal/handler"
	"github.com/danielyang21/GoBattleServer/internal/repository"
	"github.com/danielyang21/GoBattleServer/internal/service"
)

func main() {
	// Load database configuration
	dbConfig := database.LoadConfigFromEnv()

	// Create database connection pool
	pool, err := database.NewPool(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(pool)

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(pool)
	speciesRepo := repository.NewPostgresPokemonSpeciesRepository(pool)
	pokemonRepo := repository.NewPostgresUserPokemonRepository(pool)

	// Initialize services
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Initialize router
	router := handler.NewRouter(userRepo, gachaService)

	// Setup routes with middleware
	httpHandler := router.SetupRoutes()

	// Get port from environment or use default
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 API Server starting on http://localhost:%s", port)
		log.Printf("📋 Health check: http://localhost:%s/health", port)
		log.Printf("📚 API Base URL: http://localhost:%s/api", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited cleanly")
}