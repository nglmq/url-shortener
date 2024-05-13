package db

import (
	"database/sql"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	"github.com/nglmq/url-shortener/config"
	"golang.org/x/net/context"
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
 		id SERIAL PRIMARY KEY,
		alias TEXT NOT NULL,
		url TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);
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

func (s *PostgresStorage) SaveURL(ctx context.Context, alias, url string) (string, error) {
	var existingAlias string

	err := s.db.QueryRowContext(
		ctx, `
        INSERT INTO urls(alias, url) VALUES ($1, $2)
        ON CONFLICT (url) DO UPDATE SET url = EXCLUDED.url
        RETURNING alias
    `, alias, url).Scan(&existingAlias)

	if err, ok := err.(*pq.Error); ok && err.Code == pgerrcode.UniqueViolation {
		existingAlias, _ = s.GetURL(ctx, alias)

		return existingAlias, nil
	}

	return existingAlias, nil
}

//func (s *PostgresStorage) SaveBatch(urls map[string]string) error {
//	tx, err := s.db.Begin()
//	if err != nil {
//		return err
//	}
//	defer tx.Rollback()
//
//	stmt, err := tx.Prepare("INSERT INTO urls(alias, url) VALUES ($1, $2)")
//	if err != nil {
//		return err
//	}
//	defer stmt.Close()
//
//	for alias, url := range urls {
//		_, err = stmt.Exec(alias, url)
//		if err != nil {
//			return err
//		}
//	}
//
//	return tx.Commit()
//}

func (s *PostgresStorage) GetURL(ctx context.Context, alias string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT url FROM urls WHERE alias = $1", alias)

	var resURL string

	err := row.Scan(&resURL)
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
