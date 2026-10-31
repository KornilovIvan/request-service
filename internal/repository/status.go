package repository

import (
	"context"
	"database/sql"

	"request-service/internal/model"
)

type StatusRepository struct {
	db *sql.DB
}

func NewStatusRepository(db *sql.DB) *StatusRepository {
	return &StatusRepository{db: db}
}

func (r *StatusRepository) List(ctx context.Context) ([]model.Status, error) {
	const query = `
	SELECT code, title
	FROM statuses
	ORDER BY code`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []model.Status
	for rows.Next() {
		var status model.Status
		if err := rows.Scan(&status.Code, &status.Title); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	return statuses, rows.Err()
}
