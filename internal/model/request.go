package model

import "time"

type Request struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	Status     int       `json:"status"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
}
