package sql

import (
	"errors"
	"fmt"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func HandelPgErrors(err error, domain string) error {
	var pgErr *pgconn.PgError
	if errors.Is(err, pgx.ErrNoRows) {
		return &entity.NotFoundError{Field: domain, Err: err}
	}
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return &entity.AlredyExitError{Field: domain, Err: err}
		}
	}
	return fmt.Errorf("%w: %w", entity.ServiceError, err)
}
