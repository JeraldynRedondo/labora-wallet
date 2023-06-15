package service

const (
	InsertWalletInTx                   = "INSERT INTO wallets (dni_request, country_id, created_date, balance) VALUES ($1, $2, $3, $4) RETURNING *"
	UpdateWalletByID                   = "UPDATE wallets SET dni_request = $1, country_id = $2, created_date = $3, balance = $4 WHERE id = $5 RETURNING *"
	DeleteWalletByID                   = "DELETE FROM wallets WHERE id = $1"
	GetTotalWalletCount                = "SELECT COUNT(*) FROM wallets"
	GetWalletsByPage                   = "SELECT * FROM wallets ORDER BY id OFFSET $1 LIMIT $2"
	InsertLogEntry                     = "INSERT INTO logs (dni_request, country_id, status_request, date_request, request_type) VALUES ($1, $2, $3, $4, $5) RETURNING *"
	UpdateWalletBalanceAfterDeposit    = "UPDATE wallets SET balance = $1 WHERE id = $2"
	UpdateWalletBalanceAfterWithdrawal = "UPDATE wallets SET balance = $1 WHERE id = $2"
	InsertTransaction                  = "INSERT INTO transactions (wallet_id, amount, transaction_type) VALUES ($1, $2, $3) RETURNING *"
	GetWalletBalanceByID               = "SELECT balance FROM wallets WHERE id = $1"
	GetWalletByID                      = "SELECT * FROM wallets WHERE id = $1"
	UpdateBalanceValueInDeposit        = "UPDATE wallets SET balance = $1 WHERE id = $2"
	UpdateBalanceValueInWithdrawal     = "UPDATE wallets SET balance = $1 WHERE id = $2"
	GetTotalLogCount                   = "SELECT COUNT(*) FROM logs"
	GetLogsByPage                      = "SELECT * FROM logs ORDER BY id OFFSET $1 LIMIT $2"
	GetInfoWalletById                  = "SELECT id, balance FROM wallets WHERE id = $1"
	GetInfoTransWalletById             = "SELECT t.transaction_type, t.amount, t.date_transaction FROM wallets w JOIN transactions t ON w.id = t.wallet_id WHERE w.id = $1"
)
