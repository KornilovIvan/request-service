package model

import "time"

type Document struct {
	ID        int       `json:"id"`
	RequestID int       `json:"request_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
