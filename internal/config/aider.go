package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/streeboga/apihub-setup/internal/apihub"
	"github.com/streeboga/apihub-setup/internal/backup"
)

// ConfigureAider patches Aider config YAML with APIHub
func ConfigureAider(configPath, baseURL, token string, models []apihub.Model) error {
	if _, err := backup.BackupFile(configPath); err != nil {
		return fmt.Errorf("backup: %w", err)
	}

	var cfg map[string]any
	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("parsing config: %w", err)
		}
	}
	if cfg == nil {
		cfg = make(map[string]any)
	}

	cfg["openai-api-base"] = baseURL
	cfg["openai-api-key"] = token

	if len(models) > 0 {
		for _, m := range models {
			if m.Type == "chat" && len(m.Offers) > 0 {
				cfg["model"] = "openai/" + m.Offers[0].Model
				break
			}
		}
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("creating dir: %w", err)
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	if err := os.WriteFile(configPath, out, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}
