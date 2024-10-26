package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// This function reads the JSON file found at ~/.gatorconfig.json and returns a Config struct
func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	configFile := Config{}
	if err := json.Unmarshal(data, &configFile); err != nil {
		return Config{}, err
	}

	return configFile, nil
}

// This function gets path to the $HOME directory and concatenates it with '/.gatorconfig.json' 
func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path = path + "/" + configFileName
	return path, nil
}

// This method on Config struct sets the CurrentUserName field and writes the Config struct to a JSON file 
func (conifg *Config) SetUser(current_user_name string) error {
	conifg.CurrentUserName = current_user_name

	if err := write(*conifg); err != nil {
		return err
	}

	return nil
}

// This function writes/updates the config file on disk
func write(config Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0777); err != nil {
		return err
	}

	return nil
} 