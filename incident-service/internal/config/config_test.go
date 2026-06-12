package config

import "testing"

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("BOT_TOKEN", "test-token")
	t.Setenv("POSTGRES_DB", "testdb")
	t.Setenv("POSTGRES_USER", "testuser")
	t.Setenv("POSTGRES_PASSWORD", "testpass")
}

func TestLoad(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("POSTGRES_HOST", "db.example.com")
	t.Setenv("POSTGRES_PORT", "5433")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.BotToken != "test-token" {
		t.Errorf("BotToken = %q, want %q", cfg.BotToken, "test-token")
	}
	if cfg.PostgresHost != "db.example.com" {
		t.Errorf("PostgresHost = %q, want %q", cfg.PostgresHost, "db.example.com")
	}
	if cfg.PostgresPort != 5433 {
		t.Errorf("PostgresPort = %d, want %d", cfg.PostgresPort, 5433)
	}
	if cfg.PostgresDB != "testdb" {
		t.Errorf("PostgresDB = %q, want %q", cfg.PostgresDB, "testdb")
	}
	if cfg.PostgresUser != "testuser" {
		t.Errorf("PostgresUser = %q, want %q", cfg.PostgresUser, "testuser")
	}
	if cfg.PostgresPassword != "testpass" {
		t.Errorf("PostgresPassword = %q, want %q", cfg.PostgresPassword, "testpass")
	}
}

func TestLoadDefaults(t *testing.T) {
	setRequiredEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.PostgresHost != "localhost" {
		t.Errorf("default PostgresHost = %q, want %q", cfg.PostgresHost, "localhost")
	}
	if cfg.PostgresPort != 5432 {
		t.Errorf("default PostgresPort = %d, want %d", cfg.PostgresPort, 5432)
	}
}

func TestLoadMissingBotToken(t *testing.T) {
	t.Setenv("POSTGRES_DB", "testdb")
	t.Setenv("POSTGRES_USER", "testuser")
	t.Setenv("POSTGRES_PASSWORD", "testpass")

	if _, err := Load(); err == nil {
		t.Fatal("expected error when BOT_TOKEN is missing, got nil")
	}
}
