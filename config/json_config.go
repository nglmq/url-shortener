package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONConfig Struct for JSON configuration
type JSONConfig struct {
	ServerAddress   string `json:"server_address"`    // аналог переменной окружения SERVER_ADDRESS или флага -a
	BaseURL         string `json:"base_url"`          // аналог переменной окружения BASE_URL или флага -b
	FileStoragePath string `json:"file_storage_path"` // аналог переменной окружения FILE_STORAGE_PATH или флага -f
	DatabaseDSN     string `json:"database_dsn"`      // аналог переменной окружения DATABASE_DSN или флага -d
	EnableHTTPS     bool   `json:"enable_https"`      // аналог переменной окружения ENABLE_HTTPS или флага -s
	TrustedSubnet   string `json:"trusted_subnet"`    // аналог переменной окружения TRUSTED_SUBNET или флага -t
}

// ReadJSONConfig read config from file
func ReadJSONConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	var config JSONConfig

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return fmt.Errorf("error decoding json: %v", err)
	}

	if config.ServerAddress != "" {
		FlagRunAddr = config.ServerAddress
	}
	if config.BaseURL != "" {
		FlagBaseURL = config.BaseURL
	}
	if config.FileStoragePath != "" {
		FlagInMemoryStorage = config.FileStoragePath
	}
	if config.DatabaseDSN != "" {
		DBConnection = config.DatabaseDSN
	}
	if config.EnableHTTPS {
		EnableHTTPS = true
	}
	if config.TrustedSubnet != "" {
		TrustedSubnet = config.TrustedSubnet
	}

	return nil
}
