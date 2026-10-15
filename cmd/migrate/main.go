package main

import (
	"fmt"
	"log"
	"os"

	"github.com/everestp/ecom-go-api/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()

	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down>")
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseUrl,
	)
	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil {
			log.Printf("migration source close: %v", sourceErr)
		}
		if databaseErr != nil {
			log.Printf("migration database close: %v", databaseErr)
		}
	}()

	switch os.Args[1] {
	case "up":
		fmt.Println("Applying migrations...")

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration.up: %v", err)
		}

		fmt.Println("Migrations applied successfully")

	case "down":
		fmt.Println("Rolling back one migration...")

		if err := m.Steps(-1); err != nil {
			log.Fatalf("migration.down: %v", err)
		}

		fmt.Println("Migration rolled back successfully")

	default:
		log.Fatalf("unknown command %q: use up or down", os.Args[1])
	}
}
