package service

import (
	"database/sql"
	"fmt"
	"log"
	"my-labora-wallet-project/model"
	"time"
)

func (Db *PostgresDBHandler) CreateWalletInTx(wallet model.Wallet, tx *sql.Tx) (model.Wallet, error) {
	//Validation
	if wallet.DNI == "" || wallet.Country == "" {
		log.Fatal("Existen campos vacíos")
	}

	query := `INSERT INTO wallets (dni_request, country_id, date_request,balance)
	VALUES ($1, $2, $3, $4) RETURNING *`
	row := tx.QueryRow(query, &wallet.DNI, &wallet.Country, time.Now(), 100)

	err := row.Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Date_request, &wallet.Balance)
	if err != nil {
		tx.Rollback()

		return model.Wallet{}, fmt.Errorf("Error creating the wallet in the transaction: %w", err)
	}

	return wallet, nil
}

// CreateWallet is a function that creates a Wallet in the database.
func (Db *PostgresDBHandler) CreateWallet(wallet model.Wallet, log model.Log) (model.Wallet, error) {

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

	err = Db.CreateLogInTx(log, tx)
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

		query := "UPDATE wallets SET dni_request = $1, country_id = $2, date_request = $3, balance = $4 WHERE id = $5 RETURNING *"

		row := Db.QueryRow(query, &wallet.DNI, &wallet.Country, time.Now(), &wallet.Balance,id)

		err := row.Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Date_request,&wallet.Balance)

		if err != nil {


			return model.Wallet{}, fmt.Errorf("Error extracting wallet: %w", err)
		}*/

	return model.Wallet{}, nil
}

// DeleteWalletInTx it is a function that updates a wallet by id during a transaction.
func (Db *PostgresDBHandler) DeleteWalletInTx(id int, tx *sql.Tx) error {

	query := "DELETE FROM wallets WHERE id = $1"

	_, err := tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error executing delete query: %w", err)
	}

	return nil
}

// DeleteWalletInTx it is a function that updates a wallet by id during a transaction.
func (Db *PostgresDBHandler) searchWalletByIdInTx(id int, tx *sql.Tx, log *model.Log) error {
	var wallet model.Wallet
	query := "SELECT * FROM wallets WHERE id=$1"

	err := tx.QueryRow(query, id).Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Date_request, &wallet.Balance)
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error querying database for wallet: %w", err)
	}
	log.DNI = wallet.DNI
	log.Country = wallet.Country
	log.Status_request = "Deleted"
	log.Date_request = time.Now()
	log.Request_type = "DELETE WALLET"

	return nil
}

// DeleteWallet it is a function that updates a wallet by id.
func (Db *PostgresDBHandler) DeleteWallet(id int, log model.Log) error {
	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		tx.Rollback()

		return fmt.Errorf("Error at the beginning of the transaction: %w", err)
	}

	err = Db.searchWalletByIdInTx(id, tx, &log)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error trying to search the wallet in the transaction: %w", err)
	}

	err = Db.DeleteWalletInTx(id, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error trying to delete the wallet in the transaction: %w", err)
	}

	err = Db.CreateLogInTx(log, tx)
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
		err := rows.Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Date_request, &wallet.Balance)
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
func (Db *PostgresDBHandler) CreateLog(logM model.Log) error {
	//Validation
	if logM.DNI == "" || logM.Country == "" || logM.Status_request == "" || logM.Request_type == "" {
		log.Fatal("Existen campos vacíos")
	}

	// Insert the new log in the database
	query := `INSERT INTO logs (dni_request,country_id, status_request, date_request,request_type)
                        VALUES ($1, $2, $3, $4, $5) RETURNING *`
	row := Db.QueryRow(query, &logM.DNI, &logM.Country, &logM.Status_request, time.Now(), &logM.Request_type)

	err := row.Scan(&logM.DNI, &logM.Country, &logM.Status_request, time.Now(), &logM.Request_type)
	if err != nil {

		return fmt.Errorf("Error creating the log: %w", err)
	}

	return nil
}

// CreateLog is a function that creates a Log in the database during a transaction.
func (Db *PostgresDBHandler) CreateLogInTx(logM model.Log, tx *sql.Tx) error {
	//Validation
	if logM.DNI == "" || logM.Country == "" || logM.Status_request == "" || logM.Request_type == "" {
		log.Fatal("Existen campos vacíos")
	}

	// Insert the new log in the database
	query := `INSERT INTO logs (dni_request,country_id, status_request, date_request,request_type)
                        VALUES ($1, $2, $3, $4, $5) RETURNING *`
	row := tx.QueryRow(query, &logM.DNI, &logM.Country, &logM.Status_request, time.Now(), &logM.Request_type)

	err := row.Scan(&logM.ID, &logM.DNI, &logM.Country, &logM.Status_request, &logM.Date_request, &logM.Request_type)
	if err != nil {
		tx.Rollback()

		return fmt.Errorf("Error creating the log in the transaction: %w", err)
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
