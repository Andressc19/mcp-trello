package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var userHomeDir = os.UserHomeDir

type Credentials struct {
	APIKey string `json:"apiKey"`
	Token  string `json:"token"`
}

func LoadCredentials() (*Credentials, error) {
	// Priority 1: Environment variables
	envKey := os.Getenv("TRELLO_API_KEY")
	envToken := os.Getenv("TRELLO_TOKEN")

	if envKey != "" && envToken != "" {
		return &Credentials{APIKey: envKey, Token: envToken}, nil
	}

	// Priority 2: Config file
	home, err := userHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".config", "opencode", "mcp-trello.json")
		data, err := os.ReadFile(configPath)
		if err == nil {
			var creds Credentials
			if err := json.Unmarshal(data, &creds); err == nil {
				if creds.APIKey != "" && creds.Token != "" {
					return &creds, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("missing Trello credentials\n\n" +
		"Option 1 — Environment variables:\n" +
		"  export TRELLO_API_KEY=your_api_key\n" +
		"  export TRELLO_TOKEN=your_token\n\n" +
		"Option 2 — Config file at ~/.config/opencode/mcp-trello.json:\n" +
		"  {\n" +
		"    \"apiKey\": \"your_trello_api_key\",\n" +
		"    \"token\": \"your_trello_token\"\n" +
		"  }")
}
