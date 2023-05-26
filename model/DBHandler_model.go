package model

// DBHandler is an interface that implements the methods of the database.
type DBHandler interface {
	CreateWallet(wallet Wallet, log Log) (Wallet, error)
	UpdateWallet(id int, wallet Wallet) (Wallet, error)
	DeleteWallet(id int, log Log) error
	WalletStatus(pages, walletsPerPage int) ([]Wallet, int, error)
	CreateLog(log Log) error
	GetLogs(pages, logsPerPage int) ([]Log, int, error)
}
