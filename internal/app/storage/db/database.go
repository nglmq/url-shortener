package db

import (
	"database/sql"

	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	"github.com/nglmq/url-shortener/config"
	"golang.org/x/net/context"
)

// PostgresStorage is a new DB interface
type PostgresStorage struct {
	db *sql.DB
}

// InitDBConnection creates a connection and table
func InitDBConnection() (*PostgresStorage, error) {
	db, err := sql.Open("pgx", config.DBConnection)
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS urls(
 		id SERIAL PRIMARY KEY,
 		userId TEXT NOT NULL,
		alias TEXT NOT NULL,
		url TEXT NOT NULL UNIQUE,
		deleted BOOLEAN NOT NULL DEFAULT false,
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

// SaveURL saves a URL to the database
func (s *PostgresStorage) SaveURL(ctx context.Context, userID, alias, url string) (string, error) {
	var existingAlias string

	err := s.db.QueryRowContext(
		ctx, `
        INSERT INTO urls(userId, alias, url) VALUES ($1, $2, $3)
        ON CONFLICT (url) DO UPDATE SET url = EXCLUDED.url
        RETURNING alias
    `, userID, alias, url).Scan(&existingAlias)

	if err, ok := err.(*pq.Error); ok && err.Code == pgerrcode.UniqueViolation {
		existingAlias, _, _ = s.GetURL(ctx, alias)

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

// GetURL returns the URL associated with the given alias. If URL is deleted, return empty string.
func (s *PostgresStorage) GetURL(ctx context.Context, alias string) (string, bool, error) {
	row := s.db.QueryRowContext(ctx, "SELECT url, deleted FROM urls WHERE alias = $1", alias)

	var resURL string
	var deleted bool

	err := row.Scan(&resURL, &deleted)
	if err != nil {
		return "", false, err
	}

	return resURL, deleted, nil
}

// GetAllUserURLs returns all URLs associated with the given userID.
func (s *PostgresStorage) GetAllUserURLs(ctx context.Context, userID string) (map[string]string, error) {
	userURLs := make(map[string]string)

	rows, err := s.db.QueryContext(ctx, "SELECT alias, url FROM urls WHERE userId = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alias, url string
		err = rows.Scan(&alias, &url)
		if err != nil {
			return nil, err
		}

		userURLs[alias] = url
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userURLs, nil
}

// DeleteURL marks a URL as deleted
func (s *PostgresStorage) DeleteURL(ctx context.Context, alias, userID string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE urls SET deleted = true WHERE alias = $1 AND userId = $2", alias, userID)
	if err != nil {
		return err
	}

	return nil
}

// GetStats get urls and users quantity
func (s *PostgresStorage) GetStats(ctx context.Context) (int, int, error) {
	uriQuantity := 0
	usersQuantity := 0

	rowURLs := s.db.QueryRowContext(ctx, "SELECT count(DISTINCT url) FROM urls")
	if err := rowURLs.Scan(&uriQuantity); err != nil {
		return 0, 0, err
	}

	rowUsers := s.db.QueryRowContext(ctx, "SELECT count(DISTINCT userId) FROM urls")
	if err := rowUsers.Scan(&usersQuantity); err != nil {
		return 0, 0, err
	}

	return uriQuantity, usersQuantity, nil
}

// Ping the database
func (s *PostgresStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}

	return nil
}

// CloseDBConnection the database connection
func (s *PostgresStorage) CloseDBConnection() error {
	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}
