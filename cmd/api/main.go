package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/everestp/ecom-go-api/internal/config"

)

func main() {
	cfg := config.MustLoad()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	srv := http.Server{
		Addr:         ":" + cfg.Env,
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
