// Package main is the entry point for the Universal Go Backend.
//
//	@title			Universal Go Backend
//	@version		1.0
//	@description	Production-ready Go backend template.
//	@host			localhost:8080
//	@BasePath		/
//	@schemes		http https
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "my-app/docs"
	"my-app/internal/app"
	"my-app/internal/config"
	"my-app/internal/db"
	"my-app/internal/routes"
	"my-app/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := db.NewPool(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	a := app.New(cfg, pool)
	r := routes.New(cfg, a)
	srv := server.New(cfg, r)

	// Start server in background goroutine.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server starting on :%s (env=%s)", cfg.AppPort, cfg.AppEnv)
		if err := srv.Start(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until signal received.
	sig := <-quit
	log.Printf("received signal %s — initiating graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), server.ShutdownTimeout)
	defer cancel()

	log.Println("stopping HTTP server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("closing database pool...")
	pool.Close()

	log.Println("shutdown complete")
}
