package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// ConfigFile is the default configuration file used for plugin definitions.
const ConfigFile = "plugins/plugins.yaml"

type Config struct {
	ServerPort      string         `yaml:"server_port"`
	RequiredPlugins []string       `yaml:"required_plugins"`
	HelloTimeoutMs  int            `yaml:"hello_timeout_ms"`
	Plugins         []PluginConfig `yaml:"plugins"`
}

type PluginConfig struct {
	ID      string            `yaml:"id"`
	Name    string            `yaml:"name"`
	Command string            `yaml:"command"`
	Args    []string          `yaml:"args"`
	WorkDir string            `yaml:"work_dir"`
	Env     map[string]string `yaml:"env"`
	Address string            `yaml:"address"`
}

func LoadConfig(path string) (Config, error) {
	if path == "" {
		path = ConfigFile
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}
	if err != nil {
		return Config{}, fmt.Errorf("read plugin config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode plugin config: %w", err)
	}

	if cfg.ServerPort == "" {
		return Config{}, errors.New("server_port is required")
	}
	// Default hello wait timeout to 2000ms if not set or invalid.
	if cfg.HelloTimeoutMs <= 0 {
		cfg.HelloTimeoutMs = 2000
	}
	for i := range cfg.Plugins {
		if cfg.Plugins[i].ID == "" {
			cfg.Plugins[i].ID = fmt.Sprintf("plugin-%d", i+1)
		}
		if cfg.Plugins[i].Command != "" && cfg.Plugins[i].WorkDir != "" {
			if !filepath.IsAbs(cfg.Plugins[i].WorkDir) {
				cfg.Plugins[i].WorkDir = filepath.Clean(cfg.Plugins[i].WorkDir)
			}
		}
	}
	return cfg, nil
}
