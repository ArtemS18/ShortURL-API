package entity

import (
	"errors"
	"fmt"
)

var ServiceError error = errors.New("servicе error")

type NotFoundError struct {
	Err   error
	Field string
}

func (v *NotFoundError) Error() string {
	return fmt.Sprintf("not found error on field %s: %v", v.Field, v.Err)
}

func (v *NotFoundError) Unwrap() error {
	return v.Err
}

type AlredyExitError struct {
	Err   error
	Field string
}

func (v *AlredyExitError) Error() string {
	return fmt.Sprintf("already exists error on field %s: %v", v.Field, v.Err)
}

func (v *AlredyExitError) Unwrap() error {
	return v.Err
}

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

func NewValidationError(field, details string) *ValidationError {
	return &ValidationError{Err: InvalidInput, Field: field, Details: details}
}
