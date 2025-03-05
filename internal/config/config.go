package config

import (
	"encoding/json"
	// "errors"
	"fmt"
	// "log"
	"os"
	"path/filepath"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct {
	DbUrl 			string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func NewConfig() *Config {
	return &Config{}
}

func Read() (*Config, error) {
	
	file, err := openToReadConfigFile()
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}
	return config, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user

	file, err := openToWriteConfigFile()
	if err != nil {
		return fmt.Errorf("unable to set user %v: %w", user, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("unable to set user %v: %w", user, err)
	}
	
	return nil
}

func openToReadConfigFile() (*os.File, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	filepath := filepath.Join(userHomeDir, configFileName)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func openToWriteConfigFile() (*os.File, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	filepath := filepath.Join(userHomeDir, configFileName)

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}