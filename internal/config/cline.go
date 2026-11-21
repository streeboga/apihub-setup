package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/streeboga/apihub-setup/internal/apihub"
	"github.com/streeboga/apihub-setup/internal/backup"
)

// ConfigureCline patches VS Code settings for Cline extension
func ConfigureCline(configPath, baseURL, token string, models []apihub.Model) error {
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

	cfg["cline.apiProvider"] = "openai"
	cfg["cline.openAiBaseUrl"] = baseURL
	cfg["cline.openAiApiKey"] = token

	if len(models) > 0 {
		for _, m := range models {
			if m.Type == "chat" && len(m.Offers) > 0 {
				cfg["cline.openAiModelId"] = m.Offers[0].Model
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
