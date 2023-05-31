package service

import (
	"database/sql"
	"fmt"
	"log"
	"my-labora-wallet-project/model"
	"sync"
	"time"
)

func (Db *PostgresDBHandler) CreateWalletInTx(wallet model.Wallet, tx *sql.Tx) (model.Wallet, error) {
	//Validation
	if wallet.DNI == "" || wallet.Country == "" {
		log.Fatal("Existen campos vacíos")
	}

	query := `INSERT INTO wallets (dni_request, country_id, created_date,balance)
	VALUES ($1, $2, $3, $4) RETURNING *`
	row := tx.QueryRow(query, &wallet.DNI, &wallet.Country, time.Now(), 100)

	err := row.Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Created_date, &wallet.Balance)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("Error creating the wallet in the transaction: %w", err)
	}

	return wallet, nil
}

// CreateWallet is a function that creates a Wallet in the database.
func (Db *PostgresDBHandler) CreateWallet(wallet model.Wallet) (model.Wallet, error) {

	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		tx.Rollback()

		return model.Wallet{}, fmt.Errorf("Error at the beginning of the transaction: %w", err)
	}

	wallet, err = Db.CreateWalletInTx(wallet, tx)
	if err != nil {
		tx.Rollback()

		return model.Wallet{}, fmt.Errorf("Error trying to create the wallet in the transaction: %w", err)
	}

	err = Db.CreateLogInTx(wallet.DNI, wallet.Country, "Approved", "CREATE WALLET", tx)
	if err != nil {
		tx.Rollback()

		return model.Wallet{}, fmt.Errorf("Error trying to create the log in the transaction: %w", err)
	}

	// Commit the transaction if no errors occur
	err = tx.Commit()
	if err != nil {
		tx.Rollback()

		return model.Wallet{}, fmt.Errorf("Error committing the transaction: %w", err)
	}

	return wallet, nil
}

// UpdateWallet it is a function that updates a wallet by id.
func (Db *PostgresDBHandler) UpdateWallet(id int, wallet model.Wallet) (model.Wallet, error) {
	/*

		query := "UPDATE wallets SET dni_request = $1, country_id = $2, created_date = $3, balance = $4 WHERE id = $5 RETURNING *"

		row := Db.QueryRow(query, &wallet.DNI, &wallet.Country, time.Now(), &wallet.Balance,id)

		err := row.Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Created_date,&wallet.Balance)

		if err != nil {


			return model.Wallet{}, fmt.Errorf("Error extracting wallet: %w", err)
		}*/

	return model.Wallet{}, nil
}

// DeleteWallet it is a function that updates a wallet by id.
func (Db *PostgresDBHandler) DeleteWallet(id int) error {
	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		tx.Rollback()

		return fmt.Errorf("Error at the beginning of the transaction: %w", err)
	}

	wallet, err := Db.searchWalletByIdInTx(id, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error trying to search the wallet in the transaction: %w", err)
	}

	DNI, Country, status_request, request_type, err := Db.DeleteWalletInTx(wallet, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error trying to delete the wallet in the transaction: %w", err)
	}

	err = Db.CreateLogInTx(DNI, Country, status_request, request_type, tx)
	if err != nil {
		tx.Rollback()

		return fmt.Errorf("Error trying to create the log in the transaction: %w", err)
	}

	// Commit the transaction if no errors occur
	err = tx.Commit()
	if err != nil {
		tx.Rollback()

		return fmt.Errorf("Error committing the transaction: %w", err)
	}

	return nil
}

// DeleteWalletInTx it is a function that updates a wallet by id during a transaction.
func (Db *PostgresDBHandler) DeleteWalletInTx(wallet model.Wallet, tx *sql.Tx) (string, string, string, string, error) {
	query := "DELETE FROM wallets WHERE id = $1"

	_, err := tx.Exec(query, wallet.ID)
	if err != nil {
		tx.Rollback()
		return "", "", "", "", fmt.Errorf("error executing delete query: %w", err)
	}

	DNI := wallet.DNI
	Country := wallet.Country
	status_request := "Deleted"
	request_type := "DELETE WALLET"

	return DNI, Country, status_request, request_type, nil
}

// DeleteWalletInTx it is a function that updates a wallet by id during a transaction.
func (Db *PostgresDBHandler) searchWalletByIdInTx(id int, tx *sql.Tx) (model.Wallet, error) {
	var wallet model.Wallet
	query := "SELECT * FROM wallets WHERE id=$1"

	err := tx.QueryRow(query, id).Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Created_date, &wallet.Balance)
	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return model.Wallet{}, fmt.Errorf("Error querying database for wallet: %w", err)
	}

	return wallet, nil
}

// WalletStatus is a function that queries a database and returns a number of wallets per page.
func (Db *PostgresDBHandler) WalletStatus(pages, walletsPerPage int) ([]model.Wallet, int, error) {

	//Calculate the initial index and wallet limit based on the current page and wallets per page.
	start := (pages - 1) * walletsPerPage

	//Get the total number of rows in the wallets table
	var count int
	query := "SELECT COUNT(*) FROM wallets"
	err := Db.QueryRow(query).Scan(&count)
	if err != nil {

		return nil, 0, fmt.Errorf("Error querying the count in database: %w", err)
	}

	// Get the list of elements corresponding to the current page
	query = "SELECT * FROM wallets ORDER BY id OFFSET $1 LIMIT $2"
	rows, err := Db.Query(query, start, walletsPerPage)
	if err != nil {

		return nil, 0, fmt.Errorf("Error querying database: %w", err)
	}

	defer rows.Close()

	var wallets []model.Wallet

	for rows.Next() {
		var wallet model.Wallet
		err := rows.Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Created_date, &wallet.Balance)
		if err != nil {

			return nil, 0, fmt.Errorf("Error extracting wallet: %w", err)
		}
		wallets = append(wallets, wallet)
	}

	if len(wallets) == 0 {

		return nil, 0, fmt.Errorf("No wallets found for page %d", pages)
	}

	return wallets, count, nil
}

// CreateLog is a function that creates a Log in the database.
func (Db *PostgresDBHandler) CreateLog(DNI, Country, status_request, request_type string) error {
	var logM model.Log
	//Validation
	if DNI == "" || Country == "" || status_request == "" || request_type == "" {
		log.Fatal("Existen campos vacíos")
	}

	// Insert the new log in the database
	query := `INSERT INTO logs (dni_request,country_id, status_request, date_request,request_type)
                        VALUES ($1, $2, $3, $4, $5) RETURNING *`
	row := Db.QueryRow(query, DNI, Country, status_request, time.Now(), request_type)

	err := row.Scan(&logM.DNI, &logM.Country, &logM.Status_request, &logM.Date_request, &logM.Request_type)
	if err != nil {

		return fmt.Errorf("Error creating the log: %w", err)
	}

	return nil
}

// CreateLog is a function that creates a Log in the database during a transaction.
func (Db *PostgresDBHandler) CreateLogInTx(DNI, Country, status_request, request_type string, tx *sql.Tx) error {
	var logM model.Log
	//Validation
	if DNI == "" || Country == "" || status_request == "" || request_type == "" {
		log.Fatal("Existen campos vacíos")
	}

	// Insert the new log in the database
	query := `INSERT INTO logs (dni_request,country_id, status_request, date_request,request_type)
                        VALUES ($1, $2, $3, $4, $5) RETURNING *`
	row := Db.QueryRow(query, DNI, Country, status_request, time.Now(), request_type)

	err := row.Scan(&logM.DNI, &logM.Country, &logM.Status_request, &logM.Date_request, &logM.Request_type)
	if err != nil {

		return fmt.Errorf("Error creating the log: %w", err)
	}

	return nil
}

// Movement is a function that performs a money transaction from one wallet to another.
func (Db *PostgresDBHandler) CreateMovement(trans model.Transaction_Request) error {
	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		tx.Rollback()

		return fmt.Errorf("Error at the beginning of the transaction: %w", err)
	}

	validation, err := Db.validateBalanceInTx(trans.SenderID, trans.Amount, tx)
	if err != nil {
		tx.Rollback()

		return fmt.Errorf("Error trying to create the wallet in the transaction: %w", err)
	}

	if validation {

		wallet, err := Db.searchWalletByIdInTx(trans.SenderID, tx)
		if err != nil {
			tx.Rollback()

			return fmt.Errorf("Error trying to search the wallet in transaction: %w", err)
		}
		err = processTransaction(&wallet, "withdraw", trans.Amount, tx)
		if err != nil {
			tx.Rollback()

			return fmt.Errorf("Error trying to do the withdraw in transaction: %w", err)
		}
		err = Db.CreateLogInTx(wallet.DNI, wallet.Country, "Approved", "Withdraw Movement", tx)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error trying to create the log in transaction: %w", err)
		}

		err = Db.CreateTransactionInTx(wallet.ID, trans.Amount, "Deposit", tx)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error trying to create the movement in transaction: %w", err)
		}

		wallet, err = Db.searchWalletByIdInTx(trans.ReceiverID, tx)
		if err != nil {
			tx.Rollback()

			return fmt.Errorf("Error trying to search the wallet in transaction: %w", err)
		}
		err = processTransaction(&wallet, "deposit", trans.Amount, tx)
		if err != nil {
			tx.Rollback()

			return fmt.Errorf("Error trying to do the deposit in transaction: %w", err)
		}
		err = Db.CreateLogInTx(wallet.DNI, wallet.Country, "Approved", "Deposit Movement", tx)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error trying to create the log in transaction: %w", err)
		}
		err = Db.CreateTransactionInTx(wallet.ID, trans.Amount, "Deposit", tx)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error trying to create the movement in transaction: %w", err)
		}
	} else {
		wallet, err := Db.searchWalletByIdInTx(trans.SenderID, tx)
		if err != nil {
			tx.Rollback()

			return fmt.Errorf("Error trying to create the wallet in the transaction: %w", err)
		}

		err = Db.CreateLogInTx(wallet.DNI, wallet.Country, "Denied", "Transfer Movement", tx)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error trying to create the log in the transaction: %w", err)
		}
	}

	// Commit the transaction if no errors occur
	err = tx.Commit()
	if err != nil {
		tx.Rollback()

		return fmt.Errorf("Error committing the transaction: %w", err)
	}

	return nil
}

var mutex sync.Mutex

// validateBalanceInTx is a function that creates a Log in the database.
func (Db *PostgresDBHandler) validateBalanceInTx(id, amount int, tx *sql.Tx) (bool, error) {
	var balance int
	//Validation
	mutex.Lock()
	query := `SELECT balance FROM wallets WHERE id = $1`
	row := tx.QueryRow(query, id)

	err := row.Scan(balance)
	if err != nil {

		return false, fmt.Errorf("Error querying wallet balance: %w", err)
	}
	mutex.Unlock()
	if balance >= amount {

		return true, nil
	}

	return false, nil
}

// processTransaction is a function that performs a transaction (deposit or withdrawal) in the wallet
func processTransaction(wallet *model.Wallet, transactionType string, amount int, tx *sql.Tx) error {
	mutex.Lock()
	switch transactionType {
	case "deposit":
		wallet.Deposit(amount)
		query := `UPDATE wallets SET balance = $1 WHERE id = $2`
		row := tx.QueryRow(query, wallet.Balance, wallet.ID)
		err := row.Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Created_date, &wallet.Balance)
		if err != nil {

			return fmt.Errorf("Error querying wallet balance: %w", err)
		}

	case "withdraw":
		wallet.Withdraw(amount)
		query := `UPDATE wallets SET balance = $1 WHERE id = $2`
		row := tx.QueryRow(query, wallet.Balance, wallet.ID)
		err := row.Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Created_date, &wallet.Balance)
		if err != nil {

			return fmt.Errorf("Error querying wallet balance: %w", err)
		}
	}
	mutex.Unlock()

	return nil
}

func (Db *PostgresDBHandler) CreateTransactionInTx(wallet_id, amount int, transaction_type string, tx *sql.Tx) error {
	var movement model.Movement
	//Validation
	if transaction_type == "" {
		log.Fatal("Existen campos vacíos")
	}

	// Insert the new log in the database
	query := `INSERT INTO transactions (wallet_id,transaction_type,amount,date_transaction)
                        VALUES ($1, $2, $3, $4, $5) RETURNING *`
	row := tx.QueryRow(query, wallet_id, transaction_type, amount, time.Now())

	err := row.Scan(&movement.Wallet_id, &movement.Transaction_type, &movement.Amount, &movement.Date_transaction)
	if err != nil {

		return fmt.Errorf("Error creating the log: %w", err)
	}

	return nil
}

// GetLogs is a function that queries a database and returns a number of logs per page.
func (Db *PostgresDBHandler) GetLogs(pages, logsPerPage int) ([]model.Log, int, error) {
	//Calculate the initial index and log limit based on the current page and logs per page.
	start := (pages - 1) * logsPerPage

	//Get the total number of rows in the log table
	var count int
	query := "SELECT COUNT(*) FROM logs"
	err := Db.QueryRow(query).Scan(&count)
	if err != nil {

		return nil, 0, fmt.Errorf("Error querying the count in database: %w", err)
	}

	// Get the list of elements corresponding to the current page
	query = "SELECT * FROM logs ORDER BY id OFFSET $1 LIMIT $2"
	rows, err := Db.Query(query, start, logsPerPage)
	if err != nil {

		return nil, 0, fmt.Errorf("Error querying database: %w", err)
	}

	defer rows.Close()

	var logs []model.Log

	for rows.Next() {
		var log model.Log
		err := rows.Scan(&log.ID, &log.DNI, &log.Country, &log.Status_request, &log.Date_request, &log.Request_type)
		if err != nil {

			return nil, 0, fmt.Errorf("Error extracting log: %w", err)
		}
		logs = append(logs, log)
	}

	if len(logs) == 0 {

		return nil, 0, fmt.Errorf("No logs found for page %d", pages)
	}

	return logs, count, nil
}

// GetWallet is a function that queries a database for wallet id and returns the wallet with transactions.
func (Db *PostgresDBHandler) GetWalletById(id int) (model.WalletIdResponse, error) {
	return model.WalletIdResponse{}, nil
}
