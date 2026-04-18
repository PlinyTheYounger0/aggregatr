package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	err := write(*cfg)
	if err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(cfgPath)
	if err != nil {
		return Config{}, fmt.Errorf("Error Opening Config: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return Config{}, fmt.Errorf("Error Reading Config: %w", err)
	}

	cfg := Config{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("Error Decoding JSON: %w", err)
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Config File Path Error: %w", err)
	}

	filePath := filepath.Join(homePath, configFileName)
	return filePath, nil
}

func write(cfg Config) error {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Error Encoding JSON: %w", err)
	}

	if err := os.WriteFile(cfgPath, data, 0666); err != nil {
		return fmt.Errorf("Error Writing to File: %w", err)
	}

	return nil
}
