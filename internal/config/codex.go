package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/streeboga/apihub-setup/internal/apihub"
	"github.com/streeboga/apihub-setup/internal/backup"
)

// ConfigureCodex patches Codex config.toml with APIHub (text-based, preserves formatting)
func ConfigureCodex(configPath, baseURL, token string, models []apihub.Model) error {
	if _, err := backup.BackupFile(configPath); err != nil {
		return fmt.Errorf("backup: %w", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading config: %w", err)
	}

	content := string(data)

	// Set model_provider = "apihub"
	if strings.Contains(content, "model_provider") {
		// Replace existing
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "model_provider") {
				lines[i] = `model_provider = "apihub"`
			}
		}
		content = strings.Join(lines, "\n")
	} else {
		// Add at top
		content = `model_provider = "apihub"` + "\n" + content
	}

	// Add/update [model_providers.apihub] section
	section := fmt.Sprintf(`[model_providers.apihub]
name = "APIHUB"
base_url = "%s"
env_key = "APIHUB_API_KEY"
`, baseURL)

	if strings.Contains(content, "[model_providers.apihub]") {
		// Replace existing section (until next section or EOF)
		idx := strings.Index(content, "[model_providers.apihub]")
		// find next section
		nextSection := len(content)
		for i := idx + 1; i < len(content); i++ {
			if content[i] == '[' && (i == 0 || content[i-1] == '\n') {
				// check if it's a section header
				end := i
				for end < len(content) && content[end] != '\n' {
					end++
				}
				header := content[i:end]
				if strings.HasPrefix(header, "[") && strings.HasSuffix(strings.TrimSpace(header), "]") {
					nextSection = i
					break
				}
			}
		}
		content = content[:idx] + section + content[nextSection:]
	} else {
		// Append at end (before last newline if exists)
		content = strings.TrimRight(content, "\n") + "\n\n" + section
	}

	// Set default model if not present
	if !strings.Contains(content, "\nmodel = ") && !strings.HasPrefix(content, "model = ") {
		if len(models) > 0 {
			for _, m := range models {
				if m.Type == "chat" && len(m.Offers) > 0 {
					content = `model = "` + m.Offers[0].Model + `"` + "\n" + content
					break
				}
			}
		}
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("creating dir: %w", err)
	}

	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}
