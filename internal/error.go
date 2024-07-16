package internal

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	Error       string            `json:"error"`
	Validations []ValidationError `json:"validations,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

var (
	ErrVehicleAlreadyExists          = errors.New("vehicle already exists")
	ErrUniqueConstraintViolation     = errors.New("unique constraint violation")
	ErrForeignKeyConstraintViolation = errors.New("foreign key constraint violation")
	ErrNoRows                        = errors.New("no rows in result set")
	ErrInvalidParams                 = errors.New("invalid params")
)

// DBErrorToInternal converts db error to internal error
func DBErrorToInternal(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %s", ErrNoRows, err.Error())
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return fmt.Errorf("%w: %s", ErrForeignKeyConstraintViolation, pgErr.Detail)
		case "23505":
			return fmt.Errorf("%w: %s", ErrUniqueConstraintViolation, pgErr.Detail)
		default:
			return err
		}
	}
	return err
}

func RenderErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Error: msg,
	}
}
