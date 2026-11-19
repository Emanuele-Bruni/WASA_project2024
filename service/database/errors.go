package database

import (
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
)

type FailureKind int

const (
	UnknownFailure FailureKind = iota
	MissingResource
	ConflictingResource
	InvalidData
)

// ClassifyError exposes stable categories without leaking SQLite details to the API.
func ClassifyError(err error) FailureKind {
	if errors.Is(err, sql.ErrNoRows) {
		return MissingResource
	}
	var constraint sqlite3.Error
	if errors.As(err, &constraint) {
		switch constraint.ExtendedCode {
		case sqlite3.ErrConstraintUnique, sqlite3.ErrConstraintPrimaryKey:
			return ConflictingResource
		case sqlite3.ErrConstraintForeignKey:
			return MissingResource
		case sqlite3.ErrConstraintCheck, sqlite3.ErrConstraintNotNull:
			return InvalidData
		}
	}
	return UnknownFailure
}
