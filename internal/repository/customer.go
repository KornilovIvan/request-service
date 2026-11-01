package repository

import (
	"context"
	"database/sql"
	"strings"

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

func (r *CustomerRepository) Create(ctx context.Context, name string) (model.Customer, error) {
	name = strings.TrimSpace(name)

	const query = `
	INSERT INTO customers (name)
	VALUES ($1)
	RETURNING id, name, created_at`

	var customer model.Customer
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&customer.ID,
		&customer.Name,
		&customer.CreatedAt,
	)
	return customer, err
}

func (r *CustomerRepository) List(ctx context.Context) ([]model.Customer, error) {
	const query = `
	SELECT id, name, created_at
	FROM customers
	ORDER BY created_at DESC, id DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []model.Customer{}
	for rows.Next() {
		var customer model.Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.CreatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, rows.Err()
}
