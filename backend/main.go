package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Solnijko/go-knowledge-base/backend/internal/auth"
)

func createTables() error {
	dbpool := auth.ConnectDb()
	defer dbpool.Close()
	if err := auth.CreateUserTable(dbpool, context.Background()); err != nil {
		log.Fatalf("failed to create users table: %v", err)
	}
	return nil
}

func main() {

	createTables()

	mux := http.NewServeMux()
	auth.AuthRoutes(mux)

	log.Println("Server is started on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server start error: %v", err)
	}
}
