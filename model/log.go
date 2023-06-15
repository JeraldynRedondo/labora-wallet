package model

import "time"

//Log is a structure that represents the changes in the wallets that are recorded in the Logs table.
type Log struct {
	ID            int       `json:"id"`
	DNI           string    `json:"dni_request"`
	Country       string    `json:"country_id"`
	StatusRequest string    `json:"status_request"`
	DateRequest   time.Time `json:"date_request"`
	RequestType   string    `json:"request_type"`
}
