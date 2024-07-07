package database

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/glethuillier/mps/client/internal/logger"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type Database struct {
	*sql.DB
}

// CreateDatabase creates a local SQLite3 database
func CreateDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS FILES (
			ReceiptId TEXT,
			RootHash BLOB
		)
	`)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

// AddRoot adds root hashes to the database
func (db *Database) AddRootHash(ReceiptId string, RootHash []byte) error {
	query := `
		INSERT INTO FILES (ReceiptId, RootHash)
		VALUES (?, ?)
		`

	statement, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(
		ReceiptId,
		RootHash,
	)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	logger.Logger.Debug(
		"added root hash to the database",
		zap.String("receipt_id", ReceiptId),
		zap.String("root_hash", hex.EncodeToString(RootHash)),
	)

	return nil
}

// GetRootHash retrieves the root hash associated with a given
// receipt ID
func (db *Database) GetRootHash(receiptId string) (string, error) {
	var rootHash []byte
	query := `SELECT RootHash FROM FILES WHERE ReceiptId = ?`

	// Execute the query and scan the result into the rootHash variable
	err := db.QueryRow(query, receiptId).Scan(&rootHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no RootHash found for receipt ID '%s'", receiptId)
		}
		return "", err
	}

	return hex.EncodeToString(rootHash), nil
}
