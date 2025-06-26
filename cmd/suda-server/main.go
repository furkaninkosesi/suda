package main

import (
	"log"
	"net/http"
	v1 "suda-backend/internal/api/v1"
)

func main() {
	mux := http.NewServeMux()

	v1.RegisterRoutes(mux)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
