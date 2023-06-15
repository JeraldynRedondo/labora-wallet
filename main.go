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

	router.HandleFunc("/v1/wallets", controller.CreateWallet).Methods("POST")
	router.HandleFunc("/v1/wallets/{id}", controller.UpdateWallet).Methods("PUT")
	router.HandleFunc("/v1/wallets/{id}", controller.DeleteWallet).Methods("DELETE")
	router.HandleFunc("/v1/wallets", controller.WalletStatus).Methods("GET")
	router.HandleFunc("/v1/wallets/{id}", controller.GetWalletById).Methods("GET")
	router.HandleFunc("/v1/logs", controller.GetLogs).Methods("GET")
	router.HandleFunc("/v1/wallets/transaction", controller.CreateMovement).Methods("POST")

	// Configure CORS middleware
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"*"}),
	)

	// Add CORS middleware to all routes
	handler := corsOptions(router)

	if err := config.StartServer(handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
