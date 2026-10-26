package database

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Postgres error codes most often needed when mapping DB errors to API errors.
// Full list: https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
	pgCheckViolation      = "23514"
	pgNotNullViolation    = "23502"
)

// IsNoRows reports whether err means "query returned no rows".
func IsNoRows(err error) bool { return errors.Is(err, pgx.ErrNoRows) }

// IsUniqueViolation reports a UNIQUE constraint failure (e.g. duplicate email → 409).
func IsUniqueViolation(err error) bool { return hasCode(err, pgUniqueViolation) }

// IsForeignKeyViolation reports a FK failure (e.g. referenced row missing → 404/422).
func IsForeignKeyViolation(err error) bool { return hasCode(err, pgForeignKeyViolation) }

// IsCheckViolation reports a CHECK constraint failure (e.g. price <= 0 → 422).
func IsCheckViolation(err error) bool { return hasCode(err, pgCheckViolation) }

// ConstraintName returns the violated constraint's name, if any.
func ConstraintName(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.ConstraintName
	}
	return ""
}

func hasCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}
