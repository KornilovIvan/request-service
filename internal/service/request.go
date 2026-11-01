package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"request-service/internal/model"
	"request-service/internal/repository"
)

type RequestService struct {
	repo *repository.RequestRepository
}

func NewRequestService(repo *repository.RequestRepository) *RequestService {
	return &RequestService{repo: repo}
}

func (s *RequestService) Create(ctx context.Context, customerID int, title, status string) (model.Request, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return model.Request{}, ErrInvalidTitle
	}

	if err := s.repo.StatusExists(ctx, status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Request{}, ErrInvalidStatus
		}
		return model.Request{}, err
	}

	if err := s.repo.CustomerExists(ctx, customerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Request{}, ErrCustomerNotFound
		}
		return model.Request{}, err
	}

	created, err := s.repo.Create(ctx, customerID, title, status)
	if err != nil {
		return model.Request{}, err
	}
	return s.Get(ctx, created.ID)
}

func (s *RequestService) UpdateStatus(ctx context.Context, id int, status string) (model.Request, error) {
	if err := s.repo.StatusExists(ctx, status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Request{}, ErrInvalidStatus
		}
		return model.Request{}, err
	}

	if _, err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Request{}, ErrNotFound
		}
		return model.Request{}, err
	}
	return s.Get(ctx, id)
}

func (s *RequestService) List(ctx context.Context, status, createdFrom, limit *string, now time.Time) ([]model.Request, error) {
	filter, err := NormalizeListFilter(status, createdFrom, limit, now)
	if err != nil {
		return nil, err
	}
	return s.repo.List(ctx, filter.Status, filter.CreatedFrom, filter.Limit)
}

func (s *RequestService) Get(ctx context.Context, id int) (model.Request, error) {
	request, err := s.repo.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Request{}, ErrNotFound
	}
	return request, err
}
