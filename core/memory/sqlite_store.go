//go:build sqlite

package memory

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		role TEXT,
		content TEXT
	);`
	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) LogInteraction(role, content string) {
	_, err := s.db.Exec("INSERT INTO logs (role, content) VALUES (?, ?)", role, content)
	if err != nil {
		log.Printf("failed to log interaction: %v", err)
	}
}
