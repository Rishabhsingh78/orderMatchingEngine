package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Rishabhsingh78/orderMatchingEngine/internals/apis"
	"github.com/Rishabhsingh78/orderMatchingEngine/internals/engine"
)

func main() {
	// Initialize Engine
	eng := engine.NewEngine()

	// Initialize Handlers
	handler := apis.NewHandler(eng)

	// Initialize Router
	router := apis.NewRouter(handler)

	// Server Configuration
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
