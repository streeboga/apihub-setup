package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/streeboga/apihub-setup/internal/apihub"
	"github.com/streeboga/apihub-setup/internal/backup"
)

// ClaudeCodeConfig represents ~/.claude/settings.json
type ClaudeCodeConfig struct {
	Env          map[string]string `json:"env,omitempty"`
	Permissions  map[string]any    `json:"permissions,omitempty"`
	Extra        map[string]any    `json:"-"`
}

// ConfigureClaudeCode patches Claude Code settings with APIHub
func ConfigureClaudeCode(configPath, baseURL, token string, models []apihub.Model) error {
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

	// Ensure env exists
	env, ok := cfg["env"].(map[string]any)
	if !ok {
		env = make(map[string]any)
	}
	// Convert to string map
	envStr := make(map[string]string)
	for k, v := range env {
		if s, ok := v.(string); ok {
			envStr[k] = s
		}
	}

	// Set APIHub credentials
	envStr["ANTHROPIC_BASE_URL"] = baseURL
	envStr["ANTHROPIC_AUTH_TOKEN"] = token

	// Pick default models from chat models
	var haikuModel, sonnetModel, opusModel string
	for _, m := range models {
		if m.Type != "chat" {
			continue
		}
		if len(m.Offers) == 0 {
			continue
		}
		lower := m.Key
		if haikuModel == "" && (containsSubstring(lower, "haiku") || containsSubstring(lower, "mini") || containsSubstring(lower, "flash")) {
			haikuModel = m.Offers[0].Model
		}
		if sonnetModel == "" && (containsSubstring(lower, "sonnet") || containsSubstring(lower, "pro")) {
			sonnetModel = m.Offers[0].Model
		}
		if opusModel == "" && (containsSubstring(lower, "opus") || containsSubstring(lower, "ultra")) {
			opusModel = m.Offers[0].Model
		}
	}
	// Fallbacks
	if haikuModel == "" && len(models) > 0 {
		haikuModel = models[0].Offers[0].Model
	}
	if sonnetModel == "" && len(models) > 1 {
		sonnetModel = models[1].Offers[0].Model
	}
	if opusModel == "" && len(models) > 2 {
		opusModel = models[2].Offers[0].Model
	}

	if haikuModel != "" {
		envStr["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = haikuModel
		envStr["ANTHROPIC_SMALL_FAST_MODEL"] = haikuModel
	}
	if sonnetModel != "" {
		envStr["ANTHROPIC_DEFAULT_SONNET_MODEL"] = sonnetModel
	}
	if opusModel != "" {
		envStr["ANTHROPIC_DEFAULT_OPUS_MODEL"] = opusModel
	}

	// Convert back
	envAny := make(map[string]any)
	for k, v := range envStr {
		envAny[k] = v
	}
	cfg["env"] = envAny

	// Write
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
