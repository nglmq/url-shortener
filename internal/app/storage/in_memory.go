package storage

import (
	"bufio"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
)

type URLs struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func CreateFile(path string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE/os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
		return err
	}
	defer file.Close()

	log.Printf("file created %s", path)
	return nil
}

func WriteURLsToFile(path string, urls map[string]string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for shortURL, originalURL := range urls {
		id := uuid.New().String()
		url := URLs{
			UUID:        id,
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		}
		if err := encoder.Encode(url); err != nil {
			log.Printf("Error encoding URL: %v", err)
			return err
		}
	}

	return nil
}

func ReadURLsFromFile(path string, urlsMap map[string]string) error {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var url URLs
		if err := json.Unmarshal(scanner.Bytes(), &url); err != nil {
			return err
		}

		urlsMap[url.ShortURL] = url.OriginalURL
	}

	return nil
}
