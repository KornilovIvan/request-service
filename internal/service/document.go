package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"request-service/internal/model"
	"request-service/internal/repository"
)

type DocumentService struct {
	repo *repository.DocumentRepository
}

func NewDocumentService(repo *repository.DocumentRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

func (s *DocumentService) Create(ctx context.Context, requestID int, name string) (model.Document, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return model.Document{}, ErrDocumentName
	}

	document, err := s.repo.Create(ctx, requestID, name)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Document{}, ErrNotFound
	}
	return document, err
}

func (s *DocumentService) List(ctx context.Context, requestID int) ([]model.Document, error) {
	documents, err := s.repo.List(ctx, requestID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return documents, err
}

func (s *DocumentService) Get(ctx context.Context, id int) (model.Document, error) {
	document, err := s.repo.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Document{}, ErrNotFound
	}
	return document, err
}
