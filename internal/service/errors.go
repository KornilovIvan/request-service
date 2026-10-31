package service

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrDocumentName = errors.New("document name is required")
)
