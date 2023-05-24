package model

import (
	"time"
)

// Wallet is a struct that represents the Wallet object that belongs to the items table.
type Wallet struct {
	ID            int       `json:"id"`
	DNI           string    `json:"dni_request"`
	Country       string    `json:"country_id"`
	Order_request time.Time `json:"orderDate"`
}
