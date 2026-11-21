package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/streeboga/apihub-setup/internal/apihub"
	"github.com/streeboga/apihub-setup/internal/backup"
)

// ConfigureContinue patches Continue config with APIHub models
func ConfigureContinue(configPath, baseURL, token string, models []apihub.Model) error {
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

	// Build models list
	var continueModels []map[string]any
	for _, m := range models {
		if m.Type != "chat" || len(m.Offers) == 0 {
			continue
		}
		if m.Offers[0].Unavailable != nil {
			continue
		}
		continueModels = append(continueModels, map[string]any{
			"title":   m.Name,
			"provider": "openai",
			"model":   m.Offers[0].Model,
			"apiBase": baseURL,
			"apiKey":  token,
		})
	}

	cfg["models"] = continueModels

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
