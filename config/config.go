package config

import (
	"flag"
	"os"
)

var (
	FlagRunAddr         string
	FlagBaseURL         string
	FlagInMemoryStorage string
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagBaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&FlagInMemoryStorage, "f", "/tmp/short-url-db.json", "in memory storage")
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
}
