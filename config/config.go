package config

import (
	"flag"
	"os"
)

var (
	FlagRunAddr string
	FlagBaseURL string
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagBaseURL, "b", "http://localhost:8080", "base url")
	flag.Parse()

	envRunAddr := os.Getenv("SERVER_ADDRESS")
	if envRunAddr != "" {
		FlagRunAddr = envRunAddr
	} else if FlagRunAddr != "" {
		FlagRunAddr = FlagRunAddr
	} else {
		FlagRunAddr = "localhost:8080"
	}

	envBaseURL := os.Getenv("BASE_URL")
	if envBaseURL != "" {
		FlagBaseURL = envBaseURL
	} else if FlagBaseURL != "" {
		FlagBaseURL = FlagBaseURL
	} else {
		FlagBaseURL = "http://localhost:8080"
	}
}
