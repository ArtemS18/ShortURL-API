package entity

import (
	"errors"
	"fmt"
)

var ServiceError error = errors.New("servicе error")

var NotFoundError error = errors.New("slug not found error")
var AlredyExitError error = errors.New("slug already exists error")
var InvalidInput error = errors.New("invalid data")

type ValidationError struct {
	Err     error
	Field   string
	Details string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %v", v.Field, v.Err)
}

func (v *ValidationError) Unwrap() error {
	return v.Err
}

func NewValidationError(field string) *ValidationError {
	return &ValidationError{Err: InvalidInput, Field: field}
}
