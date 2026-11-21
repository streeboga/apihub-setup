package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupFile creates a backup of the given file with a timestamp suffix
func BackupFile(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // no file to backup
		}
		return "", fmt.Errorf("stat file: %w", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("path is a directory: %s", path)
	}

	backupPath := path + ".bak." + time.Now().Format("2006-01-02")
	// avoid overwriting existing backup
	if _, err := os.Stat(backupPath); err == nil {
		backupPath = path + ".bak." + time.Now().Format("2006-01-02-150405")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	if err := os.WriteFile(backupPath, data, info.Mode()); err != nil {
		return "", fmt.Errorf("writing backup: %w", err)
	}

	return backupPath, nil
}

// RestoreBackup restores a file from its backup
func RestoreBackup(backupPath, originalPath string) error {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("reading backup: %w", err)
	}
	info, err := os.Stat(backupPath)
	if err != nil {
		return fmt.Errorf("stat backup: %w", err)
	}
	return os.WriteFile(originalPath, data, info.Mode())
}

// ListBackups returns all backup files for a given path
func ListBackups(path string) ([]string, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var backups []string
	for _, e := range entries {
		if !e.IsDir() && len(e.Name()) > len(base) && e.Name()[:len(base)] == base && e.Name()[len(base):] == ".bak." {
			backups = append(backups, filepath.Join(dir, e.Name()))
		}
	}
	return backups, nil
}
