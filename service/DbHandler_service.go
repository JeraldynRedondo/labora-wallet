package service

import (
	"fmt"
	"my-labora-wallet-project/model"
	"time"
)

// CreateWallet is a function that creates a Wallet in the database.
func (Db *PostgresDBHandler) CreateWallet(wallet model.Wallet) (model.Wallet, error) {
	// Insert the new wallet in the database
	query := `INSERT INTO wallets (dni_request, country_id, order_request)
                        VALUES ($1, $2, $3) RETURNING *`
	row := Db.QueryRow(query, &wallet.DNI, &wallet.Country, time.Now())

	err := row.Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Order_request)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("Error extracting wallet: %w", err)
	}

	return wallet, nil
}

// UpdateWallet it is a function that updates a wallet by id.
func (Db *PostgresDBHandler) UpdateWallet(id int, wallet model.Wallet) (model.Wallet, error) {

	query := "UPDATE wallets SET dni_request = $1, country_id = $2, order_request = $3 WHERE id = $4 RETURNING *"
	row := Db.QueryRow(query, &wallet.DNI, &wallet.Country, time.Now(), id)
	err := row.Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Order_request)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("Error extracting wallet: %w", err)
	}

	return wallet, nil
}

// DeleteWallet it is a function that updates a wallet by id.
func (Db *PostgresDBHandler) DeleteWallet(id int) error {

	query := "DELETE FROM wallets WHERE id = ?"
	_, err := Db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error querying database: %w", err)
	}

	return nil
}

// WalletStatus it is a function that queries a database and returns a number of wallets per page.
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
		err := rows.Scan(&wallet.ID, &wallet.DNI, &wallet.Country, &wallet.Order_request)
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