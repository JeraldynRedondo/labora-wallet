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
	router.HandleFunc("/UpdateWallet/{id}", controller.UpdateWallet).Methods("PUT")
	router.HandleFunc("/DeleteWallet/{id}", controller.DeleteWallet).Methods("DELETE")
	router.HandleFunc("/WalletStatus", controller.WalletStatus).Methods("GET")
	router.HandleFunc("/wallet/{id}", controller.GetWalletById).Methods("GET")
	router.HandleFunc("/GetLogs", controller.GetLogs).Methods("GET")
	router.HandleFunc("/transaction", controller.CreateMovement).Methods("POST")

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
