package controller

import (
	"encoding/json"
	"fmt"
	"math"
	"my-labora-wallet-project/model"
	"my-labora-wallet-project/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type WalletController struct {
	WalletService service.WalletService
}

// ResponseJson is a function that sends the http response in Json format.
func ResponseJson(response http.ResponseWriter, status int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("error while marshalling object %v, trace: %+v", data, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_, err = response.Write(bytes)
	if err != nil {
		fmt.Errorf("error while writing bytes to response writer: %+v", err)
	}
}

// CreateWallet is a function that creates an Wallet from a request.
func (c *WalletController) CreateWallet(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var wallet model.Wallet

	err := json.NewDecoder(request.Body).Decode(&wallet)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Error processing the request"))
		return
	}

	wallet, err = c.WalletService.CreateWallet(wallet)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseJson(response, http.StatusOK, wallet)
}

// UpdateWallet is a function that updates an Wallet by id from a request.
func (c *WalletController) UpdateWallet(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	parameters := mux.Vars(request)
	var wallet model.Wallet

	err := json.NewDecoder(request.Body).Decode(&wallet)
	defer request.Body.Close()
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("ID must be a number"))
		return
	}

	wallet, err = c.WalletService.UpdateWallet(id, wallet)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseJson(response, http.StatusOK, wallet)
}

// DeleteWallet is a function that delete an Wallet by id from a request.
func (c *WalletController) DeleteWallet(response http.ResponseWriter, request *http.Request) {
	parameters := mux.Vars(request)
	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("ID must be a number"))
		return
	}

	err = c.WalletService.DeleteWallet(id)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseJson(response, http.StatusOK, model.Wallet{})

}

// WalletStatus is a function that returns a number of wallets per page from a request.
func (c *WalletController) WalletStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pageUser := r.URL.Query().Get("page")
	walletsUser := r.URL.Query().Get("walletsPerPage")

	page, err := strconv.Atoi(pageUser)
	if err != nil || page < 1 {
		page = 1
	}
	walletsPerPage, err := strconv.Atoi(walletsUser)
	if err != nil || walletsPerPage < 1 {
		walletsPerPage = 5
	}

	// Get the paginated list of wallets
	wallets, count, err := c.WalletService.WalletStatus(page, walletsPerPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	totalPages := int(math.Ceil(float64(count) / float64(walletsPerPage)))

	// Create a map containing information about pagination
	paginationInfo := map[string]interface{}{
		"totalPages":  totalPages,
		"currentPage": page,
	}

	// Create a map containing the list of wallets and the pagination information
	responseData := map[string]interface{}{
		"wallets":    wallets,
		"pagination": paginationInfo,
	}

	// Encode the response map in JSON format and send in the HTTP response
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
