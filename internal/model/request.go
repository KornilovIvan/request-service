package model

import "time"

type Request struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
