package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

// CreateDatabase creates a local SQLite3 database to store
// the information required to generate the proofs
func CreateDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`
		CREATE TABLE IF NOT EXISTS RECEIPTS (
			root_hash_id 	INTEGER PRIMARY KEY AUTOINCREMENT,
			receipt_id		TEXT    UNIQUE NOT NULL,
			root_hash       TEXT    UNIQUE NOT NULL
		);
		CREATE TABLE IF NOT EXISTS FILES (
			file_id      INTEGER PRIMARY KEY AUTOINCREMENT,
			root_hash_id REFERENCES RECEIPTS (root_hash_id) NOT NULL,
			filename     TEXT    NOT NULL,
			self_hash    TEXT    NOT NULL
		);
		CREATE TABLE IF NOT EXISTS TREES (
			path_id      INTEGER PRIMARY KEY AUTOINCREMENT,
			root_hash_id INTEGER REFERENCES RECEIPTS (root_hash_id),
			self_hash    TEXT    NOT NULL,
			parent_hash  TEXT    NOT NULL,
			sibling_hash TEXT    NOT NULL,
			sibling_type TEXT    NOT NULL
				CHECK (sibling_type IN ('none', 'left', 'right') ) 
		);
	`)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}
