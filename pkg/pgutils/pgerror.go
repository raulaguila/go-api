package pgutils

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrDuplicatedKey         = errors.New("duplicated key not allowed")
	ErrForeignKeyViolated    = errors.New("violates foreign key constraint")
	ErrUndefinedColumn       = errors.New("undefined column or parameter name")
	ErrDatabaseAlreadyExists = errors.New("database already exists")
)

// HandlerError handles specific PostgreSQL errors by mapping them to predefined error variables for better error management.
// It recognizes errors like duplicated keys, foreign key violations, undefined columns, and existing databases.
// If the error is not recognized as a PostgreSQL error, it prints the error and returns it unchanged.
func HandlerError(err error) error {
	if err == nil {
		return nil
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch pgError.Code {
		case "23505":
			return ErrDuplicatedKey
		case "23503":
			return ErrForeignKeyViolated
		case "42703":
			return ErrUndefinedColumn
		case "42P04":
			return ErrDatabaseAlreadyExists
		default:
			fmt.Printf("PostgreSQL error not detected: %v\n", err)
		}
	}

	return err
}
