package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/everestp/ecom-go-api/internal/config"
	"github.com/everestp/ecom-go-api/internal/db"
	"github.com/everestp/ecom-go-api/internal/handlers"
	"github.com/everestp/ecom-go-api/internal/middleware"
)

func main() {
	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect : %v", err)
	}
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	fmt.Println("database connected")
	fmt.Println("Starting ecom server")
	mux := http.NewServeMux()

	lh := handlers.NewListingHandler(db, logger)
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("POST /listings", lh.Create)

	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)
	mux.HandleFunc("POST /listings/{id}", lh.Delete)

	handler := middleware.RequestID(mux)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed %v", err)
	}

	fmt.Printf("Server is started")
}
