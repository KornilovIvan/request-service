package service

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrNameRequired     = errors.New("name is required")
	ErrDocumentName     = errors.New("document name is required")
	ErrInvalidTitle     = errors.New("invalid title")
	ErrInvalidStatus    = errors.New("invalid status")
	ErrCustomerNotFound = errors.New("customer not found")
)
