package config

import (
	"fmt"
	"net/http"
	"time"
)

// StartServer it is a function that turns on the server
func StartServer(port string, router http.Handler) error {
	servidor := &http.Server{
		Handler:      router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Starting Server on port %s...\n", port)
	if err := servidor.ListenAndServe(); err != nil {
		return fmt.Errorf("Error while starting up Server: '%v'", err)
	}
	return nil
}
