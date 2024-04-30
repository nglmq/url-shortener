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
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Директория не существует, попытка создать: %s", dir)
		if err = os.MkdirAll(dir, 0644); err != nil {
			log.Printf("Не удалось создать директорию: %v", err)
			return err
		}
		log.Println("Директория успешно создана")
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
		return err
	}
	defer file.Close()

	log.Printf("file created %s", path)
	return nil
}

func WriteURLsToFile(path string, urls map[string]string) error {
	if path == "" {
		log.Println("Path for storage is not provided, skipping file operation.")
		return nil
	}

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
	if path == "" {
		log.Println("Path for storage is not provided, skipping file operation.")
		return nil
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
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
