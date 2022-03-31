package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/portfolio-report/pr-api/server"
)

func main() {

	migrateOnly := flag.Bool("migrateOnly", false, "Migrate database and quit.")
	flag.Parse()

	cfg, db := server.PrepareApp()

	if *migrateOnly {
		os.Exit(0)
	}

	handlerConfig := server.InitializeService(cfg, db)
	router := server.CreateApp(handlerConfig)

	address := ":3000"

	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Println("Listening and serving HTTP on", address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Received signal to stop, shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shutdown gracefully: ", err)
	}

	log.Println("Server terminated.")
}
