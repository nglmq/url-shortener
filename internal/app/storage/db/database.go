package db

import (
	"database/sql"
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

	return &PostgresStorage{db: db}, nil
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
