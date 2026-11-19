package database

import (
	"database/sql"
	"net/url"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Open configures every SQLite connection, including replacements made by the pool.
// filename is a filesystem path (or :memory:), not a caller-supplied SQLite DSN.
func Open(filename string) (*sql.DB, error) {
	dsn := "file::memory:"
	if filename != ":memory:" {
		absolute, err := filepath.Abs(filename)
		if err != nil {
			return nil, err
		}
		dsn = (&url.URL{Scheme: "file", Path: absolute}).String()
	}
	return sql.Open("sqlite3", dsn+"?_foreign_keys=on&_busy_timeout=5000&_txlock=immediate")
}
