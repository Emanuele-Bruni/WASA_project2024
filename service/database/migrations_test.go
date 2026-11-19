package database

import (
	"errors"
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func migrationFixture(t *testing.T) fstest.MapFS {
	t.Helper()
	files := fstest.MapFS{}
	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		script, err := fs.ReadFile(migrationFiles, "migrations/"+entry.Name())
		if err != nil {
			t.Fatal(err)
		}
		files["migrations/"+entry.Name()] = &fstest.MapFile{Data: script}
	}
	return files
}

func TestMigrationsAreVersionedAndPreserveExistingData(t *testing.T) {
	db, c := testDatabase(t)
	id, err := db.GetOrCreateUser("alice")
	if err != nil {
		t.Fatal(err)
	}
	files := migrationFixture(t)
	files["migrations/003_next.sql"] = &fstest.MapFile{Data: []byte("CREATE TABLE migration_probe(id INTEGER PRIMARY KEY);")}
	if err = migrate(c, files); err != nil {
		t.Fatal(err)
	}
	if err = migrate(c, files); err != nil {
		t.Fatal(err)
	}
	if countRows(t, c, "schema_migrations") != 3 {
		t.Fatal("migration not applied once")
	}
	var name string
	if err = c.QueryRow("SELECT name FROM users_table WHERE id=?", id).Scan(&name); err != nil || name != "alice" {
		t.Fatal("migration lost existing data")
	}
	if err = migrate(c, migrationFiles); err == nil {
		t.Fatal("newer database version was accepted")
	}
	files["migrations/001_initial.sql"].Data = append(files["migrations/001_initial.sql"].Data, []byte("\n-- changed\n")...)
	if err = migrate(c, files); err == nil {
		t.Fatal("modified applied migration was accepted")
	}
}

func TestFailedMigrationRollsBackSchemaAndVersion(t *testing.T) {
	c, err := Open(filepath.Join(t.TempDir(), "migration.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	files := fstest.MapFS{"migrations/001_broken.sql": {Data: []byte("CREATE TABLE partial_data(id INTEGER PRIMARY KEY); INVALID SQL;")}}
	if err = migrate(c, files); err == nil {
		t.Fatal("invalid migration was accepted")
	}
	var count int
	if err = c.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("failed initial migration left partial schema or history")
	}
	if _, err = New(c); err != nil {
		t.Fatalf("clean initialization after rollback: %v", err)
	}
	files = migrationFixture(t)
	files["migrations/003_broken.sql"] = &fstest.MapFile{Data: []byte("CREATE TABLE partial_data(id INTEGER PRIMARY KEY); INVALID SQL;")}
	if err = migrate(c, files); err == nil {
		t.Fatal("invalid incremental migration accepted")
	}
	if countRows(t, c, "schema_migrations") != 2 {
		t.Fatal("failed migration advanced version")
	}
	if err = c.QueryRow("SELECT count(*) FROM sqlite_master WHERE name='partial_data'").Scan(&count); err != nil || count != 0 {
		t.Fatal("partial table survived rollback")
	}
}

func TestLegacyDatabaseIsRefusedWithoutModification(t *testing.T) {
	c, err := Open(filepath.Join(t.TempDir(), "legacy.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if _, err = c.Exec("CREATE TABLE users_table(id INTEGER PRIMARY KEY,name TEXT,conversation_ids TEXT); INSERT INTO users_table VALUES(1,'alice','2,3');"); err != nil {
		t.Fatal(err)
	}
	if _, err = New(c); !errors.Is(err, ErrLegacySchema) {
		t.Fatalf("unexpected legacy result: %v", err)
	}
	var list string
	if err = c.QueryRow("SELECT conversation_ids FROM users_table WHERE id=1").Scan(&list); err != nil || list != "2,3" {
		t.Fatal("legacy data changed")
	}
	var count int
	if err = c.QueryRow("SELECT count(*) FROM sqlite_master WHERE name='schema_migrations'").Scan(&count); err != nil || count != 0 {
		t.Fatal("legacy schema changed")
	}
}
