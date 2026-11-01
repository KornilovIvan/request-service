package service

import (
	"errors"
	"testing"
	"time"
)

func TestNormalizeListFilter(t *testing.T) {
	today := time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)

	str := func(value string) *string { return &value }
	date := func(year int, month time.Month, day int) *time.Time {
		value := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		return &value
	}

	tests := []struct {
		name          string
		status        *string
		createdFrom   *string
		limit         *string
		wantStatusSet bool
		wantStatus    string
		wantCreated   *time.Time
		wantLimit     int
		wantErr       error
	}{
		{
			name:      "no status",
			wantLimit: 20,
		},
		{
			name:          "empty status",
			status:        str(""),
			wantStatusSet: true,
			wantStatus:    "",
			wantLimit:     20,
		},
		{
			name:          "active",
			status:        str("active"),
			wantStatusSet: true,
			wantStatus:    "active",
			wantLimit:     20,
		},
		{
			name:    "done",
			status:  str("done"),
			wantErr: ErrInvalidStatus,
		},
		{
			name:      "no created_from",
			wantLimit: 20,
		},
		{
			name:        "past date",
			createdFrom: str("2022-01-01"),
			wantCreated: date(2022, time.January, 1),
			wantLimit:   20,
		},
		{
			name:        "today",
			createdFrom: str("2026-10-31"),
			wantCreated: date(2026, time.October, 31),
			wantLimit:   20,
		},
		{
			name:        "tomorrow",
			createdFrom: str("2026-11-01"),
			wantErr:     ErrInvalidDate,
		},
		{
			name:        "impossible date",
			createdFrom: str("2022-13-40"),
			wantErr:     ErrInvalidDate,
		},
		{
			name:      "no limit",
			wantLimit: 20,
		},
		{
			name:      "limit 20",
			limit:     str("20"),
			wantLimit: 20,
		},
		{
			name:      "limit 1",
			limit:     str("1"),
			wantLimit: 1,
		},
		{
			name:      "limit 100",
			limit:     str("100"),
			wantLimit: 100,
		},
		{
			name:    "limit 0",
			limit:   str("0"),
			wantErr: ErrInvalidLimit,
		},
		{
			name:    "limit 101",
			limit:   str("101"),
			wantErr: ErrInvalidLimit,
		},
		{
			name:    "limit abc",
			limit:   str("abc"),
			wantErr: ErrInvalidLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeListFilter(tt.status, tt.createdFrom, tt.limit, today)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantStatusSet {
				if got.Status == nil {
					t.Fatal("status filter disabled, want enabled")
				}
				if *got.Status != tt.wantStatus {
					t.Fatalf("status = %q, want %q", *got.Status, tt.wantStatus)
				}
			} else if got.Status != nil {
				t.Fatalf("status filter enabled with %q, want disabled", *got.Status)
			}
			if tt.wantCreated == nil {
				if got.CreatedFrom != nil {
					t.Fatalf("createdFrom = %v, want nil", got.CreatedFrom)
				}
			} else if got.CreatedFrom == nil || !got.CreatedFrom.Equal(*tt.wantCreated) {
				t.Fatalf("createdFrom = %v, want %v", got.CreatedFrom, *tt.wantCreated)
			}
			if got.Limit != tt.wantLimit {
				t.Fatalf("limit = %d, want %d", got.Limit, tt.wantLimit)
			}
		})
	}
}
