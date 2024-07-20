package config

import (
	"flag"
	"os"
)

// Flags for the server config
var (
	FlagRunAddr         string
	FlagBaseURL         string
	FlagInMemoryStorage string
	DBConnection        string
)

// ParseFlags parses the command line args and ENV variables
func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagBaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&FlagInMemoryStorage, "f", "", "in memory storage")
	flag.StringVar(&DBConnection, "d", "", "postgres connection url")
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
}
