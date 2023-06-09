package service

import "my-labora-wallet-project/model"

type WalletService struct {
	DbHandler model.DBHandler
}

// CreateWallet implements the function CreateWallet of the database in a DBHandler.
func (s *WalletService) CreateWallet(wallet model.Wallet) (model.Wallet, error) {
	return s.DbHandler.CreateWallet(wallet)
}

// UpdateWallet implements the function UpdateWallet of the database in a DBHandler.
func (s *WalletService) UpdateWallet(id int, wallet model.Wallet) (model.Wallet, error) {
	return s.DbHandler.UpdateWallet(id, wallet)
}

// DeleteWallet implements the function DeleteWallet of the database in a DBHandler.
func (s *WalletService) DeleteWallet(id int) error {
	return s.DbHandler.DeleteWallet(id)
}

// WalletStatus implements the function WalletStatus of the database in a DBHandler.
func (s *WalletService) WalletStatus(pages, walletsPerPage int) ([]model.Wallet, int, error) {
	return s.DbHandler.WalletStatus(pages, walletsPerPage)
}

// CreateLog implements the function CreateLog of the database in a DBHandler.
func (s *WalletService) CreateLog(DNI, Country, status_request, request_type string) error {
	return s.DbHandler.CreateLog(DNI, Country, status_request, request_type)
}

// GetLogs implements the function GetLogs of the database in a DBHandler.
func (s *WalletService) GetLogs(pages, logsPerPage int) ([]model.Log, int, error) {
	return s.DbHandler.GetLogs(pages, logsPerPage)
}

// GetWalletById implements the function GetWalletById of the database in a DBHandler.
func (s *WalletService) GetWalletById(id int) (model.WalletIdResponse, error) {
	return s.DbHandler.GetWalletById(id)
}

// CreateMovement implements the function CreateMovement of the database in a DBHandler.
func (s *WalletService) CreateMovement(trans model.Transaction_Request) (string, error) {
	return s.DbHandler.CreateMovement(trans)
}
