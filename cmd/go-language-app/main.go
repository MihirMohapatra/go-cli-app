package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MihirMohapatra/go-cli-app/internal/todos"
)

func main() {
	port := flag.Int("port", 8080, "HTTP server port")
	flag.Parse()

	databaseURL := getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/go_cli_app?sslmode=disable")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := todos.OpenPostgres(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connect to postgres: %v", err)
	}
	defer db.Close()

	if err := todos.Migrate(ctx, db); err != nil {
		log.Fatalf("migrate postgres: %v", err)
	}

	store := todos.NewPostgresStore(db)
	handler := todos.NewHandler(store)
	addr := fmt.Sprintf(":%d", *port)

	log.Printf("listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
