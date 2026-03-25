package main

import (
	"os"
	"path/filepath"
	"testing"
)

func setupConfigTest(t *testing.T) string {
	t.Helper()
	// Isolate from real home dir
	tmpDir := t.TempDir()
	mockUserHomeDir(t, tmpDir)
	// Also set HOME env var for compatibility with other code that may read it
	t.Setenv("HOME", tmpDir)
	// Clear env vars
	t.Setenv("TRELLO_API_KEY", "")
	t.Setenv("TRELLO_TOKEN", "")
	return tmpDir
}

func writeConfigFile(t *testing.T, content string) {
	t.Helper()
	home := os.Getenv("HOME")
	dir := filepath.Join(home, ".config", "opencode")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("creating config dir: %v", err)
	}
	path := filepath.Join(dir, "mcp-trello.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing config file: %v", err)
	}
}

func TestLoadCredentials_EnvVars(t *testing.T) {
	_ = setupConfigTest(t)
	t.Setenv("TRELLO_API_KEY", "env-key")
	t.Setenv("TRELLO_TOKEN", "env-token")

	creds, err := LoadCredentials()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.APIKey != "env-key" {
		t.Errorf("expected APIKey=env-key, got %s", creds.APIKey)
	}
	if creds.Token != "env-token" {
		t.Errorf("expected Token=env-token, got %s", creds.Token)
	}
}

func TestLoadCredentials_ConfigFile(t *testing.T) {
	_ = setupConfigTest(t)
	writeConfigFile(t, `{"apiKey":"file-key","token":"file-token"}`)

	creds, err := LoadCredentials()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.APIKey != "file-key" {
		t.Errorf("expected APIKey=file-key, got %s", creds.APIKey)
	}
	if creds.Token != "file-token" {
		t.Errorf("expected Token=file-token, got %s", creds.Token)
	}
}

func TestLoadCredentials_EnvVarsPriority(t *testing.T) {
	_ = setupConfigTest(t)
	t.Setenv("TRELLO_API_KEY", "env-key")
	t.Setenv("TRELLO_TOKEN", "env-token")
	writeConfigFile(t, `{"apiKey":"file-key","token":"file-token"}`)

	creds, err := LoadCredentials()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.APIKey != "env-key" {
		t.Errorf("expected env vars to take priority, got APIKey=%s", creds.APIKey)
	}
	if creds.Token != "env-token" {
		t.Errorf("expected env vars to take priority, got Token=%s", creds.Token)
	}
}

func TestLoadCredentials_PartialEnvVars(t *testing.T) {
	_ = setupConfigTest(t)
	t.Setenv("TRELLO_API_KEY", "env-key")
	// TRELLO_TOKEN not set
	writeConfigFile(t, `{"apiKey":"file-key","token":"file-token"}`)

	creds, err := LoadCredentials()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.APIKey != "file-key" {
		t.Errorf("expected fallback to file, got APIKey=%s", creds.APIKey)
	}
	if creds.Token != "file-token" {
		t.Errorf("expected fallback to file, got Token=%s", creds.Token)
	}
}

func TestLoadCredentials_MissingCredentials(t *testing.T) {
	_ = setupConfigTest(t)
	// No env vars, no config file

	_, err := LoadCredentials()
	if err == nil {
		t.Fatal("expected error for missing credentials")
	}
}

func TestLoadCredentials_InvalidJSON(t *testing.T) {
	_ = setupConfigTest(t)
	writeConfigFile(t, `not valid json {{{`)

	_, err := LoadCredentials()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
