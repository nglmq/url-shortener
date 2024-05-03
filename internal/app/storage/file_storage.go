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
	log.Printf("create path: %s", path)

	dir := filepath.Dir(path)
	log.Printf("dir: %s", dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Директория не существует, попытка создать: %s", dir)
		if err = os.MkdirAll(dir, 0644); err != nil {
			log.Printf("Не удалось создать директорию: %v", err)
			return err
		}
		log.Println("Директория успешно создана")
		log.Println(dir)
		log.Println(filepath.Abs(dir))
	}

	//file, err := os.Create(path) //os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	//if err != nil {
	//	log.Fatalf("Failed to create file: %v", err)
	//	return err
	//}
	//defer file.Close()

	return nil
}

func WriteURLsToFile(path, alias, url string) error {
	log.Printf("write path: %s", path)
	if path == "" {
		log.Println("Path for storage is not provided, skipping file operation.")
		return nil
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	log.Println(filepath.Abs(filepath.Dir(path)))

	encoder := json.NewEncoder(file)

	id := uuid.New().String()
	newData := URLs{
		UUID:        id,
		ShortURL:    alias,
		OriginalURL: url,
	}
	if err := encoder.Encode(newData); err != nil {
		log.Printf("Error encoding URL: %v", err)
		return err
	}

	return nil
}

func ReadURLsFromFile(path string, urlsMap map[string]string) error {
	log.Printf("read path: %s", path)
	if path == "" {
		log.Println("Path for storage is not provided, skipping file operation.")
		return nil
	}

	dir := filepath.Dir(path)
	log.Println(filepath.Abs(dir))

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
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
	log.Println("URLs read from file successfully")
	return nil
}
