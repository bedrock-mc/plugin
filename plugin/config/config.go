package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	ID      string   `yaml:"id"`
	Name    string   `yaml:"name"`
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
	WorkDir struct {
		RemoteGit bool   `yaml:"remote_git"`
		Path      string `yaml:"path"`
	} `yaml:"work_dir"`
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
		pl := &cfg.Plugins[i]
		if pl.ID == "" {
			pl.ID = fmt.Sprintf("plugin-%d", i+1)
		}
		if pl.Command != "" && pl.WorkDir.Path != "" {
			if pl.WorkDir.RemoteGit {
				path := "/tmp/" + pl.ID
				err := os.RemoveAll(path)
				if err != nil {
					return cfg, err
				}

				cmd := exec.Command("git", "clone", pl.WorkDir.Path, path, "--depth=1")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err = cmd.Run()
				if err != nil {
					return cfg, err
				}

				pl.WorkDir.Path = path
			}

			if !filepath.IsAbs(pl.WorkDir.Path) {
				pl.WorkDir.Path = filepath.Clean(pl.WorkDir.Path)
			}
		}
	}
	return cfg, nil
}
