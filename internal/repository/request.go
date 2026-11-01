package repository

import (
	"context"
	"database/sql"

	"request-service/internal/model"
)

type RequestRepository struct {
	db *sql.DB
}

func NewRequestRepository(db *sql.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) Create(ctx context.Context, customerID int, title, status string) (model.Request, error) {
	const query = `
	INSERT INTO requests (customer_id, title, status)
	VALUES ($1, $2, $3)
	RETURNING id, customer_id, title, status, created_at`

	var request model.Request
	err := r.db.QueryRowContext(ctx, query, customerID, title, status).Scan(
		&request.ID,
		&request.CustomerID,
		&request.Title,
		&request.Status,
		&request.CreatedAt,
	)
	return request, err
}

func (r *RequestRepository) Get(ctx context.Context, id int) (model.Request, error) {
	const query = `
	SELECT id, customer_id, title, status, created_at
	FROM requests
	WHERE id = $1`

	var request model.Request
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.CustomerID,
		&request.Title,
		&request.Status,
		&request.CreatedAt,
	)
	return request, err
}

func (r *RequestRepository) StatusExists(ctx context.Context, code string) error {
	const query = `
	SELECT code
	FROM statuses
	WHERE code = $1`

	var found string
	return r.db.QueryRowContext(ctx, query, code).Scan(&found)
}

func (r *RequestRepository) CustomerExists(ctx context.Context, id int) error {
	const query = `
	SELECT id
	FROM customers
	WHERE id = $1`

	var found int
	return r.db.QueryRowContext(ctx, query, id).Scan(&found)
}
