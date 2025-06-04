package main

import (
	"log"
	"net/http"

	"github.com/Solnijko/go-knowledge-base/cmd/server/backend/auth"
)

func main() {

	mux := http.NewServeMux()
	auth.AuthRoutes(mux)

	log.Println("Server is started on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server start error: %v", err)
	}
}
