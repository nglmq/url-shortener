package config

import (
	"flag"
	"log/slog"
	"os"
)

// server_address": "localhost:8080", // аналог переменной окружения SERVER_ADDRESS или флага -a
// "base_url": "http://localhost", // аналог переменной окружения BASE_URL или флага -b
// "file_storage_path": "/path/to/file.db", // аналог переменной окружения FILE_STORAGE_PATH или флага -f
// "database_dsn": "", // аналог переменной окружения DATABASE_DSN или флага -d
// "enable_https": true // аналог переменной окружения ENABLE_HTTPS или флага -s
// Flags for the server config
var (
	FlagRunAddr         string
	FlagBaseURL         string
	FlagInMemoryStorage string
	DBConnection        string
	EnableHTTPS         bool
	ReadConfigFile      string
	TrustedSubnet       string
)

// ParseFlags parses the command line args and ENV variables
func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagBaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&FlagInMemoryStorage, "f", "", "in memory storage")
	flag.StringVar(&DBConnection, "d", "", "postgres connection url")
	flag.BoolVar(&EnableHTTPS, "s", false, "enable https")
	flag.StringVar(&ReadConfigFile, "c", "", "read json config from file")
	flag.StringVar(&TrustedSubnet, "t", "", "CIDR")
	flag.Parse()

	envRunAddr := os.Getenv("SERVER_ADDRESS")
	if envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}

	envBaseURL := os.Getenv("BASE_URL")
	if envBaseURL != "" {
		FlagBaseURL = envBaseURL
	}

	envInMemoryStorage := os.Getenv("FILE_STORAGE_PATH")
	if envInMemoryStorage != "" {
		FlagInMemoryStorage = envInMemoryStorage
	}

	envDBConnection := os.Getenv("DATABASE_DSN")
	if envDBConnection != "" {
		DBConnection = envDBConnection
	}

	envEnableHTTPS := os.Getenv("ENABLE_HTTPS")
	if envEnableHTTPS != "" {
		EnableHTTPS = true
	}

	envTrustedSubnet := os.Getenv("TRUSTED_SUBNET")
	if envTrustedSubnet != "" {
		TrustedSubnet = envTrustedSubnet
	}

	envReadConfig := os.Getenv("CONFIG")
	if envReadConfig != "" {
		ReadConfigFile = envReadConfig
	}

	if ReadConfigFile != "" {
		err := ReadJSONConfig(ReadConfigFile)
		if err != nil {
			slog.Info("error while reading json")
		}
	}
}
