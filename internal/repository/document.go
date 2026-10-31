package repository

import (
	"context"
	"database/sql"

	"request-service/internal/model"
)

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(ctx context.Context, requestID int, name string) (model.Document, error) {
	if err := r.lookupRequest(ctx, requestID); err != nil {
		return model.Document{}, err
	}

	const query = `
	INSERT INTO documents (request_id, name)
	VALUES ($1, $2)
	RETURNING id, request_id, name, created_at`

	var document model.Document
	err := r.db.QueryRowContext(ctx, query, requestID, name).Scan(
		&document.ID,
		&document.RequestID,
		&document.Name,
		&document.CreatedAt,
	)
	return document, err
}

func (r *DocumentRepository) List(ctx context.Context, requestID int) ([]model.Document, error) {
	if err := r.lookupRequest(ctx, requestID); err != nil {
		return nil, err
	}

	const query = `
	SELECT id, request_id, name, created_at
	FROM documents
	WHERE request_id = $1
	ORDER BY created_at DESC, id DESC`

	rows, err := r.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := []model.Document{}
	for rows.Next() {
		var document model.Document
		if err := rows.Scan(
			&document.ID,
			&document.RequestID,
			&document.Name,
			&document.CreatedAt,
		); err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	return documents, rows.Err()
}

func (r *DocumentRepository) Get(ctx context.Context, id int) (model.Document, error) {
	const query = `
	SELECT id, request_id, name, created_at
	FROM documents
	WHERE id = $1`

	var document model.Document
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&document.ID,
		&document.RequestID,
		&document.Name,
		&document.CreatedAt,
	)
	return document, err
}

func (r *DocumentRepository) lookupRequest(ctx context.Context, requestID int) error {
	const query = `
	SELECT id
	FROM requests
	WHERE id = $1`

	var id int
	return r.db.QueryRowContext(ctx, query, requestID).Scan(&id)
}
