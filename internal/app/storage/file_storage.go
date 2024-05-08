package storage

import (
	"bufio"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"os"
	"sync"
)

type FileStorage struct {
	file   *os.File
	writer *bufio.Writer
	mx     sync.RWMutex
}

type URLs struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewFileStorage(filename string) (*FileStorage, error) {
	//dir := filepath.Dir(filename)
	//if _, err := os.Stat(dir); os.IsNotExist(err) {
	//	os.MkdirAll(dir, 0755)
	//}

	log.Printf("Initializing File Storage with filename: %s", filename)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return nil, err
	}
	writer := bufio.NewWriter(file)
	return &FileStorage{file: file, writer: writer}, nil
}

func (f *FileStorage) WriteURLsToFile(alias, url string) error {
	f.mx.Lock()
	defer f.mx.Unlock()

	log.Printf("Writing to file: alias=%s, url=%s", alias, url)
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

	log.Printf("Data written and flushed successfully")
	return nil
}

func (f *FileStorage) Close() error {
	log.Println("Closing file storage")
	return f.file.Close()
}

func (f *FileStorage) ReadURLsFromFile(urlsMap map[string]string) error {
	scanner := bufio.NewScanner(f.file)

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
