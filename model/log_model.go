package model

import "time"

type Log struct {
	ID             int       `json:"id"`
	DNI            string    `json:"dni_request"`
	Status_request string    `json:"status_request"`
	Order_request  time.Time `json:"order_request"`
}
