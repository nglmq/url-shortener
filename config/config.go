package config

import (
	"flag"
	"github.com/nglmq/url-shortener/internal/app/storage"
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
	} else if FlagRunAddr != "" {
		return
	} else {
		FlagRunAddr = "localhost:8080"
	}

	envBaseURL := os.Getenv("BASE_URL")
	if envBaseURL != "" {
		FlagBaseURL = envBaseURL
	} else if FlagBaseURL != "" {
		return
	} else {
		FlagBaseURL = "http://localhost:8080"
	}

	envInMemoryStorage := os.Getenv("FILE_STORAGE_PATH")
	if envInMemoryStorage != "" {
		FlagInMemoryStorage = envInMemoryStorage
		storage.CreateNewFile(FlagInMemoryStorage)
	} else if FlagInMemoryStorage != "" {
		storage.CreateNewFile(FlagInMemoryStorage)
		return
	} else {
		FlagInMemoryStorage = "/tmp/short-url-db.json"
		storage.CreateNewFile(FlagInMemoryStorage)
	}
}
