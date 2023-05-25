package controller

import (
	"encoding/json"
	"fmt"
	"math"
	"my-labora-wallet-project/model"
	"my-labora-wallet-project/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type WalletController struct {
	WalletService service.WalletService
}

// ResponseJson it is a function that sends the http response in Json format.
func ResponseJson(response http.ResponseWriter, status int, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("error while marshalling object %v, trace: %+v", data, err)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_, err = response.Write(bytes)
	if err != nil {
		return fmt.Errorf("error while writing bytes to response writer: %+v", err)
	}

	return nil
}

// CreateWallet is a function that creates an Wallet from a request.
func (c *WalletController) CreateWallet(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var Body_request model.API_Request
	var wallet model.Wallet
	var log model.Log

	err := json.NewDecoder(request.Body).Decode(&Body_request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	wallet.DNI = Body_request.National_id
	wallet.Country = Body_request.Country
	wallet.Order_request = time.Now()
	autorization, err := service.GetApproval(Body_request.National_id, Body_request.Country, Body_request.Typ, Body_request.UserAuthorized)

	if !autorization {
		log.DNI = wallet.DNI
		log.Status_request = "Denied"
		log.Order_request = time.Now()
		c.WalletService.CreateLog(log)
		ResponseJson(response, http.StatusConflict, nil)

		return
	}

	wallet, err = c.WalletService.CreateWallet(wallet)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)

		return
	}

	log.DNI = wallet.DNI
	log.Status_request = "Approved"
	log.Order_request = time.Now()
	c.WalletService.CreateLog(log)

	ResponseJson(response, http.StatusOK, wallet)
}

// UpdateWallet is a function that updates an Wallet by id from a request.
func (c *WalletController) UpdateWallet(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	parameters := mux.Vars(request)
	var wallet model.Wallet

	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("ID must be a number"))
		return
	}

	err = json.NewDecoder(request.Body).Decode(&wallet)
	defer request.Body.Close()
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
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

// GetLogs is a function that returns a number of logs per page from a request.
func (c *WalletController) GetLogs(w http.ResponseWriter, r *http.Request) {
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

	// Get the paginated list of logs
	wallets, count, err := c.WalletService.GetLogs(page, walletsPerPage)
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
