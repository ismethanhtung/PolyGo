// @title PolyGo API
// @version 1.0
// @description High-performance Polymarket API proxy with caching and WebSocket support
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@polygo.dev

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name POLY-API-KEY

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/polygo/internal/api"
	"github.com/polygo/internal/cache"
	"github.com/polygo/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create cache
	c, err := cache.New(&cfg.Cache)
	if err != nil {
		log.Fatalf("Failed to create cache: %v", err)
	}

	// Create and start server
	server, err := api.NewServer(cfg, c)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Println("Shutting down server...")
		if err := server.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 PolyGo server starting on %s", addr)
	log.Printf("📚 Swagger UI: http://%s/swagger/index.html", addr)
	
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
