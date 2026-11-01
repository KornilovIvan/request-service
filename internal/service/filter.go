package service

import (
	"strconv"
	"time"
)

type ListFilter struct {
	Status      *string
	CreatedFrom *time.Time
	Limit       int
}

func NormalizeListFilter(status, createdFrom, limit *string, now time.Time) (ListFilter, error) {
	filter := ListFilter{Limit: 20}

	if status != nil && *status != "" {
		switch *status {
		case "active", "pending", "closed":
			value := *status
			filter.Status = &value
		default:
			return ListFilter{}, ErrInvalidStatus
		}
	}

	if createdFrom != nil {
		parsed, err := time.Parse("2006-01-02", *createdFrom)
		if err != nil {
			return ListFilter{}, ErrInvalidDate
		}
		nowUTC := now.UTC()
		today := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)
		if parsed.After(today) {
			return ListFilter{}, ErrInvalidDate
		}
		filter.CreatedFrom = &parsed
	}

	if limit != nil {
		n, err := strconv.Atoi(*limit)
		if err != nil || n < 1 || n > 100 {
			return ListFilter{}, ErrInvalidLimit
		}
		filter.Limit = n
	}

	return filter, nil
}
