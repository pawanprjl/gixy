package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Config struct {
	Profiles       map[string]Profile `json:"profiles"`
	PathMappings   map[string]string  `json:"path_mappings,omitempty"`
	DefaultProfile string             `json:"default_profile,omitempty"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".config", "gixy", "config"), nil
}

func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{Profiles: make(map[string]Profile)}, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}
	if cfg.PathMappings == nil {
		cfg.PathMappings = make(map[string]string)
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	// Atomic write: temp file + rename, so a crash mid-write can't truncate the config.
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "config-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp config: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return fmt.Errorf("chmod temp config: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp config: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp config: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replace config: %w", err)
	}
	return nil
}
