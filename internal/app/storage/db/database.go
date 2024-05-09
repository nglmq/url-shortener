package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nglmq/url-shortener/config"
)

type PostgresStorage struct {
	db *sql.DB
}

func InitDBConnection() (*PostgresStorage, error) {
	db, err := sql.Open("pgx", config.DBConnection)
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS urls(
-- 		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL,
		url TEXT NOT NULL);
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) SaveURL(alias, url string) error {
	stmt, err := s.db.Prepare("INSERT INTO urls(alias, url) VALUES ($1, $2)")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(alias, url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *PostgresStorage) GetURL(alias string) (string, error) {
	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias = $1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		return "", err
	}

	return resURL, nil
}

func (s *PostgresStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) CloseDBConnection() error {
	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}
