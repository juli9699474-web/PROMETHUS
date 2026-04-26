package router

import "errors"

import "github.com/jackc/pgx/v5"

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
