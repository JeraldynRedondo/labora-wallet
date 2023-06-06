package model

import "time"

//Log is a structure that represents the changes in the wallets that are recorded in the Logs table.
type Log struct {
	ID             int       `json:"id"`
	DNI            string    `json:"dni_request"`
	Country        string    `json:"country_id"`
	Status_request string    `json:"status_request"`
	Date_request   time.Time `json:"date_request"`
	Request_type   string    `json:"request_type"`
}
