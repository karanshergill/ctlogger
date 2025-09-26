package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sql.DB
}

type DomainEntry struct {
	ID        int64     `json:"id"`
	Domain    string    `json:"domain"`
	Timestamp time.Time `json:"timestamp"`
}

func NewDatabase(dbPath string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	db := &Database{conn: conn}
	if err := db.createTables(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return db, nil
}

func (db *Database) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS domains (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain TEXT NOT NULL UNIQUE,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_domain ON domains(domain);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON domains(timestamp);
	`

	_, err := db.conn.Exec(query)
	return err
}

func (db *Database) InsertDomain(domain string) error {
	query := `INSERT OR IGNORE INTO domains (domain) VALUES (?)`
	_, err := db.conn.Exec(query, domain)
	return err
}

func (db *Database) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}