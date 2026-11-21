package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/streeboga/apihub-setup/internal/apihub"
	"github.com/streeboga/apihub-setup/internal/backup"
)

// OpenCodeModel represents a model entry in OpenCode config
type OpenCodeModel struct {
	Name       string            `json:"name"`
	ToolCall   bool              `json:"tool_call,omitempty"`
	Reasoning  bool              `json:"reasoning,omitempty"`
	Attachment bool              `json:"attachment,omitempty"`
	Limit      *OpenCodeLimit    `json:"limit,omitempty"`
	Variants   map[string]any    `json:"variants,omitempty"`
}

type OpenCodeLimit struct {
	Context int `json:"context"`
	Output  int `json:"output"`
}

type OpenCodeProvider struct {
	Name    string                     `json:"name"`
	NPM     string                     `json:"npm"`
	Options OpenCodeProviderOptions    `json:"options"`
	Models  map[string]OpenCodeModel   `json:"models"`
}

type OpenCodeProviderOptions struct {
	BaseURL string `json:"baseURL"`
	APIKey  string `json:"apiKey"`
}

type OpenCodeConfig struct {
	Schema    string                       `json:"$schema,omitempty"`
	MCP       map[string]any               `json:"mcp,omitempty"`
	Provider  map[string]OpenCodeProvider  `json:"provider,omitempty"`
	Model     string                       `json:"model,omitempty"`
	Agent     map[string]any               `json:"agent,omitempty"`
}

// ConfigureOpenCode patches the OpenCode config with APIHub models
func ConfigureOpenCode(configPath, baseURL, token string, models []apihub.Model) error {
	// Backup existing config
	if _, err := backup.BackupFile(configPath); err != nil {
		return fmt.Errorf("backup: %w", err)
	}

	// Read existing config
	var cfg OpenCodeConfig
	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("parsing config: %w", err)
		}
	}

	if cfg.Provider == nil {
		cfg.Provider = make(map[string]OpenCodeProvider)
	}

	// Build models map
	modelsMap := make(map[string]OpenCodeModel)
	for _, m := range models {
		if m.Type != "chat" {
			continue
		}
		// Use first available offer
		if len(m.Offers) == 0 {
			continue
		}
		offer := m.Offers[0]
		// Skip unavailable
		if offer.Unavailable != nil {
			continue
		}

		ocModel := OpenCodeModel{
			Name: m.Name,
		}

		// Context length
		ctx := 128000 // default
		if m.ContextLength != nil {
			ctx = *m.ContextLength
		} else if offer.Model != "" {
			// try to infer from known models
		}

		ocModel.Limit = &OpenCodeLimit{
			Context: ctx,
			Output:  32000,
		}

		// Detect capabilities from model name/maker
		lowerName := m.Name
		lowerKey := m.Key
		_ = lowerName
		_ = lowerKey

		// Reasoning models
		if containsAny(m.Key, []string{"deepseek", "o1", "o3", "claude", "thinking", "reason", "gemini-2.5", "gemini-3"}) ||
			containsAny(m.Name, []string{"DeepSeek", "Claude", "o1", "o3", "Gemini 2.5", "Gemini 3"}) {
			ocModel.Reasoning = true
			ocModel.Variants = map[string]any{
				"low":    map[string]any{},
				"medium": map[string]any{},
				"high":   map[string]any{},
			}
		}

		// Attachment models (vision)
		if containsAny(m.Key, []string{"vision", "gpt-4", "claude", "gemini"}) ||
			containsAny(m.Name, []string{"GPT-4", "Claude", "Gemini", "Vision"}) {
			ocModel.Attachment = true
		}

		// Tool call - assume all chat models support it
		ocModel.ToolCall = true

		modelsMap[offer.Model] = ocModel
	}

	cfg.Provider["apihub"] = OpenCodeProvider{
		Name: "APIHUB",
		NPM:  "@ai-sdk/openai-compatible",
		Options: OpenCodeProviderOptions{
			BaseURL: baseURL,
			APIKey:  token,
		},
		Models: modelsMap,
	}

	// Set default model if not set
	if cfg.Model == "" && len(modelsMap) > 0 {
		// pick first model
		for k := range modelsMap {
			cfg.Model = "apihub/" + k
			break
		}
	}

	// Write config
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

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if contains(s, sub) {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
