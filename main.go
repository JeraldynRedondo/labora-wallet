package main

import (
	"log"
	"my_api_project/config"
	"my_api_project/controller"
	"my_api_project/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	var dbHandler, error = service.Connect_DB()
	if error != nil {
		log.Fatal(error)
	}
	itemService := &service.ItemService{DbHandler: dbHandler}
	controller := &controller.ItemController{ItemService: *itemService}
	router := mux.NewRouter()

	router.HandleFunc("/items", controller.GetAllItems).Methods("GET")
	router.HandleFunc("/items/page", controller.GetItemsPaginated).Methods("GET")
	//router.HandleFunc("/items/details/{id}", controller.ItemDetails).Methods("GET")
	router.HandleFunc("/items/id/{id}", controller.GetItemById).Methods("GET")
	router.HandleFunc("/items/name/{name}", controller.GetItemByName).Methods("GET")

	router.HandleFunc("/items", controller.CreateItem).Methods("POST")
	router.HandleFunc("/items/{id}", controller.UpdateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", controller.DeleteItem).Methods("DELETE")

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
