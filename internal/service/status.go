package service

import (
	"context"

	"request-service/internal/model"
	"request-service/internal/repository"
)

type StatusService struct {
	repo *repository.StatusRepository
}

func NewStatusService(repo *repository.StatusRepository) *StatusService {
	return &StatusService{repo: repo}
}

func (s *StatusService) List(ctx context.Context) ([]model.Status, error) {
	return s.repo.List(ctx)
}
