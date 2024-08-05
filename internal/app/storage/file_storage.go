package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

// FileStorage is a struct for file storage
type FileStorage struct {
	file   *os.File
	writer *bufio.Writer
	mx     sync.RWMutex
}

// URLs JSON format for saving URLs to file
type URLs struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// NewFileStorage creates a new file
func NewFileStorage(filename string) (*FileStorage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return nil, err
	}

	writer := bufio.NewWriter(file)
	return &FileStorage{file: file, writer: writer}, nil
}

// WriteURLsToFile writes the given URLs to the file
func (f *FileStorage) WriteURLsToFile(alias, url string) error {
	f.mx.Lock()
	defer f.mx.Unlock()

	urls := URLs{
		UUID:        uuid.New().String(),
		ShortURL:    alias,
		OriginalURL: url,
	}

	data, err := json.Marshal(urls)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return err
	}

	if _, err := f.writer.Write(data); err != nil {
		log.Printf("Error writing data to file: %v", err)
		return err
	}

	if err := f.writer.WriteByte('\n'); err != nil {
		log.Printf("Error writing newline to file: %v", err)
		return err
	}

	if err := f.writer.Flush(); err != nil {
		log.Printf("Error flushing writer: %v", err)
		return err
	}

	return nil
}

// Close closes the file
func (f *FileStorage) Close() error {
	return f.file.Close()
}

// ReadURLsFromFile reads the URLs to map
func (f *FileStorage) ReadURLsFromFile(urlsMap map[string]string) error {
	scanner := bufio.NewScanner(f.file)

	for scanner.Scan() {
		var url URLs
		if err := json.Unmarshal(scanner.Bytes(), &url); err != nil {
			return err
		}

		urlsMap[url.ShortURL] = url.OriginalURL
	}

	return nil
}
