package errors

import "errors"

var (
	ErrNotFound            = errors.New("record not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInternalServerError = errors.New("internal server error")
)
