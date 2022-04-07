// Package errors provides error definitions.
package errors

import (
	"errors"
	"fmt"
	"strings"
)

type Response struct {
	Status  int
	Type    string
	Message []string
}

var (
	ErrInvalidTimestamp = errors.New("Invalid timestamps: start < end")
	ErrOldTimestamp     = errors.New("Invalid timestamps: start should be > current time")
	ErrBBDDScan         = errors.New("scan bbdd failed")
	ErrQueryContext     = errors.New("query context Failed")
	ErrPrepareContext   = errors.New("prepare context failed")
	ErrExecContext      = errors.New("exec Context failed")
	ErrRowsAffected     = errors.New("could not get rows affected")
	ErrBeginTx          = errors.New("begin failed")
	ErrIDRequired       = errors.New("id is required to use the API")
	ErrFailedToInsert   = errors.New("failed to insert")
	ErrNotFound         = errors.New("not found")
	ErrTraderOpenMarket = errors.New("traders should only be able to have one market open at a time")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInputNotFound    = errors.New("input setting not found")
	ErrNotUpdated       = errors.New("data was not updated")
)

func CleanPQError(err error) error {
	if err == nil {
		return nil
	}
	return errors.New(strings.TrimPrefix(err.Error(), "pq: "))
}

// RowsAffectedError returns a formatted affected rows error.
func RowsAffectedError(n int64) error {
	return fmt.Errorf("psql: expected 1 row affected, got %d", n)
}
