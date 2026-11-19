package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConfigurationDefaultsAndEnvironment(t *testing.T) {
	t.Setenv("CFG_CONFIG_PATH", filepath.Join(t.TempDir(), "absent.yml"))
	t.Setenv("CFG_DB_FILENAME", "data/wasatext-v2.db")
	t.Setenv("CFG_WEB_API_HOST", "0.0.0.0:3000")
	cfg, err := loadConfigurationFrom(nil)
	if err != nil || cfg.DB.Filename != "data/wasatext-v2.db" {
		t.Fatalf("default configuration: %+v, %v", cfg, err)
	}
	t.Setenv("CFG_DB_FILENAME", "/data/custom.db")
	t.Setenv("CFG_WEB_API_HOST", "127.0.0.1:3100")
	t.Setenv("CFG_WEB_READ_TIMEOUT", "9s")
	cfg, err = loadConfigurationFrom(nil)
	if err != nil || cfg.DB.Filename != "/data/custom.db" || cfg.Web.APIHost != "127.0.0.1:3100" || cfg.Web.ReadTimeout != 9*time.Second {
		t.Fatalf("environment configuration: %+v, %v", cfg, err)
	}
	cfg, err = loadConfigurationFrom([]string{"--db-filename", "/data/flag.db"})
	if err != nil || cfg.DB.Filename != "/data/flag.db" {
		t.Fatalf("CLI override: %+v, %v", cfg, err)
	}
}

func TestConfigurationFileOverridesEnvironment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte("db:\n  filename: /data/file.db\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CFG_CONFIG_PATH", path)
	t.Setenv("CFG_DB_FILENAME", "/data/environment.db")
	cfg, err := loadConfigurationFrom(nil)
	if err != nil || cfg.DB.Filename != "/data/file.db" {
		t.Fatalf("YAML override: %+v, %v", cfg, err)
	}
}
