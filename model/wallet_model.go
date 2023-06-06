package model

import (
	"time"
)

// Wallet is a struct that represents the Wallet object that belongs to the wallets table.
type Wallet struct {
	ID           int       `json:"id"`
	DNI          string    `json:"dni_request"`
	Country      string    `json:"country_id"`
	Created_date time.Time `json:"date_request"`
	Balance      int       `json:"balance"`
}

// Transaction_Request is a structure that represents the request body in a transaction.
type Transaction_Request struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
	Amount     int `json:"amount"`
}

// Movement is a struct that represents the movement in a Wallet object.
type Movement struct {
	ID               int       `json:"id"`
	Wallet_id        int       `json:"wallet_id"`
	Transaction_type string    `json:"transaction_type"`
	Amount           int       `json:"amount"`
	Date_transaction time.Time `json:"date_transaction"`
}

// WalletIdResponse is a structure that represents the query response by Id of a wallet and has attributes of a wallet type object and a transaction type object.
type WalletIdResponse struct {
	ID        int            `json:"id"`
	Balance   int            `json:"balance"`
	Movements []MovementById `json:"movements"`
}

// MovementById is a struct that represents the movement in a Wallet object when searching by id.
type MovementById struct {
	Transaction_type string    `json:"transaction_type"`
	Amount           int       `json:"amount"`
	Date_transaction time.Time `json:"date_transaction"`
}

// Deposit is a method that increases the wallet balance by the given amount
func (wallet *Wallet) Deposit(amount int) {
	wallet.Balance += amount
}

// Withdraw is a Method that decreases the wallet balance with the given amount
func (wallet *Wallet) Withdraw(amount int) {
	if wallet.Balance >= amount {
		wallet.Balance -= amount
	}
}
