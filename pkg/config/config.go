package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	Token string `json:"token"`
}

var configFileName = "config.json"

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// logrus.Debugln("config dir:", configDir)

	// CLI-specific directory
	cliDir := filepath.Join(configDir, ".dblyser")
	if err := os.MkdirAll(cliDir, os.ModePerm); err != nil {
		return "", err
	}

	return filepath.Join(cliDir, configFileName), nil
}

// Save saves the config to the file
func (c *Config) Save() error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	return encoder.Encode(c)
}

// Load loads the config from the file
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil // Return empty config if not found
		}
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
