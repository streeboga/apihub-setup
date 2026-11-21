package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Tool struct {
	Name        string
	ConfigPath  string
	Exists      bool
}

// DetectTools scans the system for supported AI coding tools
func DetectTools() []Tool {
	var tools []Tool

	home, _ := os.UserHomeDir()

	// OpenCode
	opencodePath := filepath.Join(home, ".config", "opencode", "opencode.json")
	tools = append(tools, Tool{
		Name:       "opencode",
		ConfigPath: opencodePath,
		Exists:     fileExists(opencodePath),
	})

	// Claude Code
	claudePath := filepath.Join(home, ".claude", "settings.json")
	tools = append(tools, Tool{
		Name:       "claude-code",
		ConfigPath: claudePath,
		Exists:     fileExists(claudePath),
	})

	// Codex (OpenAI)
	codexPath := filepath.Join(home, ".codex", "config.toml")
	tools = append(tools, Tool{
		Name:       "codex",
		ConfigPath: codexPath,
		Exists:     fileExists(codexPath),
	})

	// Cursor
	var cursorPath string
	if runtime.GOOS == "darwin" {
		cursorPath = filepath.Join(home, "Library", "Application Support", "Cursor", "User", "settings.json")
	} else if runtime.GOOS == "windows" {
		cursorPath = filepath.Join(home, "AppData", "Roaming", "Cursor", "User", "settings.json")
	} else {
		cursorPath = filepath.Join(home, ".config", "Cursor", "User", "settings.json")
	}
	tools = append(tools, Tool{
		Name:       "cursor",
		ConfigPath: cursorPath,
		Exists:     fileExists(cursorPath),
	})

	// Continue
	continuePath := filepath.Join(home, ".continue", "config.json")
	tools = append(tools, Tool{
		Name:       "continue",
		ConfigPath: continuePath,
		Exists:     fileExists(continuePath),
	})

	// Aider
	aiderPath := filepath.Join(home, ".aider.conf.yml")
	tools = append(tools, Tool{
		Name:       "aider",
		ConfigPath: aiderPath,
		Exists:     fileExists(aiderPath),
	})

	// Cline (VS Code extension settings)
	var clinePath string
	if runtime.GOOS == "darwin" {
		clinePath = filepath.Join(home, "Library", "Application Support", "Code", "User", "settings.json")
	} else if runtime.GOOS == "windows" {
		clinePath = filepath.Join(home, "AppData", "Roaming", "Code", "User", "settings.json")
	} else {
		clinePath = filepath.Join(home, ".config", "Code", "User", "settings.json")
	}
	tools = append(tools, Tool{
		Name:       "cline",
		ConfigPath: clinePath,
		Exists:     fileExists(clinePath),
	})

	return tools
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (t Tool) String() string {
	status := "not found"
	if t.Exists {
		status = "found"
	}
	return fmt.Sprintf("%s: %s (%s)", t.Name, t.ConfigPath, status)
}
