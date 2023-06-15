package config

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// StartServer it is a function that turns on the server
func StartServer(router http.Handler) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file: %w", err)
	}

	port := os.Getenv("PORT_NUMBER")
	if port == "" {
		port = "8000"
	}

	server := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Starting Server on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("Error while starting up Server: '%v'", err)
	}

	return nil
}
