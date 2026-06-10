package sql

import (
	"errors"
	"fmt"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func HandelPgErrors(err error) error {
	var pgErr *pgconn.PgError
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %w", entity.NotFoundError, err)
	}
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("%w: %w", entity.AlredyExitError, err)
		}
	}
	return fmt.Errorf("%w: %w", entity.ServiceError, err)
}
