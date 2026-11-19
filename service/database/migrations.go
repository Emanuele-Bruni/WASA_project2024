package database

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"strings"
	"time"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

var ErrLegacySchema = errors.New("database non versionato: usare un nuovo file per V2; il database precedente non è stato modificato")

func migrate(c *sql.DB, files fs.FS) error {
	entries, err := fs.ReadDir(files, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	tx, err := c.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var versioned bool
	if err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type='table' AND name='schema_migrations')").Scan(&versioned); err != nil {
		return err
	}
	if !versioned {
		var count int
		if err = tx.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&count); err != nil {
			return err
		}
		if count != 0 {
			return ErrLegacySchema
		}
	}
	if _, err = tx.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
 version INTEGER NOT NULL PRIMARY KEY,
 name TEXT NOT NULL UNIQUE,
 checksum TEXT NOT NULL,
 applied_at TEXT NOT NULL
 )`); err != nil {
		return err
	}
	type appliedMigration struct{ name, checksum string }
	applied := make(map[int]appliedMigration)
	rows, err := tx.Query("SELECT version,name,checksum FROM schema_migrations ORDER BY version")
	if err != nil {
		return err
	}
	for rows.Next() {
		var version int
		var m appliedMigration
		if err = rows.Scan(&version, &m.name, &m.checksum); err != nil {
			rows.Close()
			return err
		}
		if version != len(applied)+1 {
			rows.Close()
			return errors.New("non-contiguous migration history")
		}
		applied[version] = m
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	if len(entries) == 0 || len(applied) > len(entries) {
		return errors.New("database version is newer than the supported schema or migrations are missing")
	}
	for i, entry := range entries {
		version, err := strconv.Atoi(strings.SplitN(entry.Name(), "_", 2)[0])
		if err != nil || entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") || version != i+1 {
			return fmt.Errorf("invalid migration sequence: %s", entry.Name())
		}
		script, err := fs.ReadFile(files, "migrations/"+entry.Name())
		if err != nil {
			return err
		}
		checksum := fmt.Sprintf("%x", sha256.Sum256(script))
		if previous, ok := applied[version]; ok {
			if previous.name != entry.Name() || previous.checksum != checksum {
				return fmt.Errorf("applied migration %d has changed", version)
			}
			continue
		}
		if _, err = tx.Exec(string(script)); err != nil {
			return fmt.Errorf("apply migration %s: %w", entry.Name(), err)
		}
		if entry.Name() == "002_messages_and_media.sql" {
			if err = migrateLegacyMedia(tx); err != nil {
				return fmt.Errorf("convert legacy media: %w", err)
			}
		}
		if _, err = tx.Exec("INSERT INTO schema_migrations(version,name,checksum,applied_at) VALUES (?,?,?,?)", version, entry.Name(), checksum, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit migrations: %w", err)
	}
	return nil
}
