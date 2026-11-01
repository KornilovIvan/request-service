package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

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

func (s *CustomerService) Create(ctx context.Context, name string) (model.Customer, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return model.Customer{}, ErrNameRequired
	}
	return s.repo.Create(ctx, name)
}

func (s *CustomerService) List(ctx context.Context) ([]model.Customer, error) {
	return s.repo.List(ctx)
}
