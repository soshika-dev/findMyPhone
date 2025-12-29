package repository

import (
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// isDuplicateError reports whether the provided error originates from a duplicate key violation.
func isDuplicateError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
