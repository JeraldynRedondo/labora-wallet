package main

import (
	"log"
	"my-labora-wallet-project/config"
	"my-labora-wallet-project/controller"
	"my-labora-wallet-project/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	var dbHandler, error = service.Connect_DB()
	if error != nil {
		log.Fatal(error)
	}
	walletService := &service.WalletService{DbHandler: dbHandler}
	controller := &controller.WalletController{WalletService: *walletService}

	router := mux.NewRouter()

	router.HandleFunc("/CreateWallet", controller.CreateWallet).Methods("POST")
	router.HandleFunc("/UpdateWallet", controller.UpdateWallet).Methods("PUT")
	router.HandleFunc("/DeleteWallet", controller.DeleteWallet).Methods("DELETE")
	router.HandleFunc("/WalletStatus", controller.WalletStatus).Methods("GET")

	// Configure CORS middleware
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST"}),
	)

	// Add CORS middleware to all routes
	handler := corsOptions(router)

	portNumber := ":9999"
	if err := config.StartServer(portNumber, handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
