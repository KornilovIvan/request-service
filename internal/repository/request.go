package repository

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

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

func (r *RequestRepository) List(ctx context.Context, status *string, createdFrom *time.Time, limit int) ([]model.Request, error) {
	query := `
	SELECT r.id, r.customer_id, r.title, r.status, r.created_at,
	       c.name AS customer_name
	FROM requests r
	JOIN customers c ON c.id = r.customer_id`
	args := make([]any, 0, 3)
	conditions := make([]string, 0, 2)
	if status != nil {
		args = append(args, *status)
		conditions = append(conditions, "r.status = $"+strconv.Itoa(len(args)))
	}
	if createdFrom != nil {
		args = append(args, *createdFrom)
		conditions = append(conditions, "r.created_at >= $"+strconv.Itoa(len(args)))
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	args = append(args, limit)
	query += " ORDER BY r.created_at DESC LIMIT $" + strconv.Itoa(len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []model.Request{}
	for rows.Next() {
		var request model.Request
		if err := rows.Scan(
			&request.ID,
			&request.CustomerID,
			&request.Title,
			&request.Status,
			&request.CreatedAt,
			&request.CustomerName,
		); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, rows.Err()
}

func (r *RequestRepository) Get(ctx context.Context, id int) (model.Request, error) {
	const query = `
	SELECT r.id, r.customer_id, r.title, r.status, r.created_at,
	       c.name AS customer_name
	FROM requests r
	JOIN customers c ON c.id = r.customer_id
	WHERE r.id = $1`

	var request model.Request
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.CustomerID,
		&request.Title,
		&request.Status,
		&request.CreatedAt,
		&request.CustomerName,
	)
	return request, err
}

func (r *RequestRepository) UpdateStatus(ctx context.Context, id int, status string) (model.Request, error) {
	const query = `
	UPDATE requests
	SET status = $1
	WHERE id = $2
	RETURNING id, customer_id, title, status, created_at`

	rows, err := r.db.QueryContext(ctx, query, status, id)
	if err != nil {
		return model.Request{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return model.Request{}, err
		}
		return model.Request{}, sql.ErrNoRows
	}

	var request model.Request
	if err := rows.Scan(
		&request.ID,
		&request.CustomerID,
		&request.Title,
		&request.Status,
		&request.CreatedAt,
	); err != nil {
		return model.Request{}, err
	}
	return request, rows.Err()
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
