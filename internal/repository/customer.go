package repository

import (
	"context"
	"database/sql"

	"request-service/internal/model"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Get(ctx context.Context, id int) (model.Customer, error) {
	const query = `
	SELECT id, name, created_at
	FROM customers
	WHERE id = $1`

	var customer model.Customer
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID,
		&customer.Name,
		&customer.CreatedAt,
	)
	return customer, err
}
