package main

import (
	"log"
	"my-labora-wallet-project/config"
	"my-labora-wallet-project/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	var _, error = service.Connect_DB()
	if error != nil {
		log.Fatal(error)
	}
	router := mux.NewRouter()

	//router.HandleFunc("/CreateWallet", ).Methods("GET")

	// Configure CORS middleware
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST"}),
	)

	// Add CORS middleware to all routes
	handler := corsOptions(router)

	portNumber := ":3000"
	if err := config.StartServer(portNumber, handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
