package model

// DBHandler is an interface that implements the methods of the database.
type DBHandler interface {
	CreateWallet(wallet Wallet) (Wallet, error)
	UpdateWallet(id int, wallet Wallet) (Wallet, error)
	DeleteWallet(id int) error
	WalletStatus(pages, walletsPerPage int) ([]Wallet, int, error)
	GetWalletById(id int) (WalletIdResponse, error)
	CreateLog(DNI, Country, status_request, request_type string) error
	GetLogs(pages, logsPerPage int) ([]Log, int, error)
	CreateMovement(trans Transaction_Request) (string, error)
}
