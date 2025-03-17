package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	DbURL string `json:"db_url"`
	User  string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(homeDir, configFileName)
	return configPath, nil
}

func ReadConfig() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func WriteConfig(c *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	config, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer config.Close()

	_, err = config.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func SetUser(c *Config) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	c.User = currentUser.Username
	return nil
}
