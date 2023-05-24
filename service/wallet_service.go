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
