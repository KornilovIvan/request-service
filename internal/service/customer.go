package service

import (
	"context"
	"database/sql"
	"errors"

	"request-service/internal/model"
	"request-service/internal/repository"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Get(ctx context.Context, id int) (model.Customer, error) {
	customer, err := s.repo.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Customer{}, ErrNotFound
	}
	return customer, err
}
