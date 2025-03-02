package config

import (
	"encoding/json"
	"os"
)

const configFileName = "/.gatorconfig.json"

func Read() (Config, error) {
	cfgFile, err := GetConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	cfgJson, err := os.ReadFile(cfgFile)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = json.Unmarshal(cfgJson, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}

func GetConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homePath + configFileName, nil
}

func write(cfg Config) error {
	cfgFile, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return nil
	}

	return os.WriteFile(cfgFile, jsonData, os.ModePerm)
}
