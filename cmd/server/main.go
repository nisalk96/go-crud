package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"restapi/internal/config"
	"restapi/internal/handlers"
	"restapi/internal/repository"
	"restapi/internal/router"
	"restapi/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		log.Fatalf("upload dir: %v", err)
	}

	ctx := context.Background()
	client, err := repository.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(shutdownCtx); err != nil {
			log.Printf("mongo disconnect: %v", err)
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping: %v", err)
	}

	db := client.Database(cfg.MongoDatabase)
	movieRepo := repository.NewMovieRepository(db, cfg.MongoCollectionMovies)
	coverStore := &storage.CoverStorage{Dir: cfg.UploadDir, MaxBytes: cfg.MaxUploadBytes}
	movieHandler := &handlers.MovieHandler{
		Repo:            movieRepo,
		Covers:          coverStore,
		PublicCoverPath: "/api/v1/files/covers",
	}
	h := router.New(movieHandler, cfg.UploadDir)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
