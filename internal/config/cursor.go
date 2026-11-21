package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/streeboga/apihub-setup/internal/apihub"
	"github.com/streeboga/apihub-setup/internal/backup"
)

// ConfigureCursor patches Cursor settings.json with APIHub
func ConfigureCursor(configPath, baseURL, token string, models []apihub.Model) error {
	if _, err := backup.BackupFile(configPath); err != nil {
		return fmt.Errorf("backup: %w", err)
	}

	var cfg map[string]any
	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("parsing config: %w", err)
		}
	}
	if cfg == nil {
		cfg = make(map[string]any)
	}

	cfg["openai.apiKey"] = token
	cfg["openai.apiBase"] = baseURL

	// Set default model
	if len(models) > 0 {
		for _, m := range models {
			if m.Type == "chat" && len(m.Offers) > 0 {
				cfg["openai.model"] = m.Offers[0].Model
				break
			}
		}
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("creating dir: %w", err)
	}
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	if err := os.WriteFile(configPath, out, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}
