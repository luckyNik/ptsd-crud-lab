package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5432/db")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	repo := NewGuitarRepo(dbpool)
	guitarHandler := NewGuitarHandler(repo)

	r := chi.NewRouter()
	guitarHandler.RegisterRoutes(r)

	http.ListenAndServe(":8080", r)
}
