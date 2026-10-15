package main

import (
	"fmt"
	"log"
	"os"

	"github.com/everestp/ecom-go-api/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(os.Args)
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down>")
	}
	m, err := migrate.New(
		"files://migrations",
		cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("migration.new %v", err)

	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatalf("migration.up: %s", err)
		}
	case "down":
		if err := m.Down(); err != nil {
			log.Fatalf("migration.down: %s", err)
		}
	default:
		log.Fatalf("unknown command %s", os.Args[1])
	}
	fmt.Println("Running migration")
}
