package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/everestp/ecom-go-api/internal/config"
	"github.com/everestp/ecom-go-api/internal/db"
	"github.com/everestp/ecom-go-api/internal/handlers"
)

func main() {
	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect : %v", err)
	}
	fmt.Println("database connected")
	fmt.Println("Starting ecom server")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /listings", handlers.List(db))

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed %v", err)
	}

	fmt.Printf("Server is started")
}
