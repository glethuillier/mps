package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/glethuillier/fvs/lib/pkg/proofs"
	"github.com/glethuillier/fvs/server/internal/common"
	"github.com/glethuillier/fvs/server/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SaveTree saves a receipt ID and the corresponding Merkle tree in the database
func (db *Database) SaveTree(receiptId uuid.UUID, tree *common.Tree) error {
	var err error
	if err = db.addRootHash(receiptId, tree.RootHash); err != nil {
		return err
	}

	for k, v := range tree.FilenameToHash {
		if err = db.addFile(tree.RootHash, k, v); err != nil {
			return err
		}
	}

	for k, v := range tree.Nodes {
		if err = db.addSubtree(
			tree.RootHash,
			k,
			v.Parent,
			v.Sibling,
			v.SiblingType,
		); err != nil {
			return err
		}
	}

	return nil
}

// addRootHash saves a root hash corresponding to a given receipt ID
// in the database
func (db *Database) addRootHash(receiptId uuid.UUID, rootHash string) error {
	query := `
	INSERT INTO RECEIPTS (receipt_id, root_hash)
	VALUES (?, ?)
	`

	statement, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(
		receiptId.String(),
		rootHash,
	)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	logger.Logger.Debug(
		"added root hash to the database",
		zap.String("root_hash", rootHash),
	)

	return nil
}

// addFile saves a filename and the corresponding file hash
// in the database
func (db *Database) addFile(rootHash, filename, selfHash string) error {
	// get the root hash ID corresponding to the root hash
	var rootHashID int
	query := "SELECT root_hash_id FROM RECEIPTS WHERE root_hash = ?"
	err := db.QueryRow(query, rootHash).Scan(&rootHashID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("root_hash %s not found", rootHash)
		}
		return err
	}

	// insert the filename and its hash
	insertQuery := "INSERT INTO FILES (root_hash_id, filename, self_hash) VALUES (?, ?, ?)"
	_, err = db.Exec(insertQuery, rootHashID, filename, selfHash)
	if err != nil {
		return err
	}

	return nil
}

// addSubtree adds a subtree (self, sibling, parent) into the database
func (db *Database) addSubtree(
	rootHash, self, parent, sibling string,
	siblingType proofs.SiblingType,
) error {
	// get the root hash ID corresponding to the root hash
	var rootHashID int
	query := "SELECT root_hash_id FROM RECEIPTS WHERE root_hash = ?"
	err := db.QueryRow(query, rootHash).Scan(&rootHashID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("root_hash %s not found", rootHash)
		}
		return err
	}

	// insert the subtree
	query = `
		INSERT INTO TREES (root_hash_id, self_hash, parent_hash, sibling_hash, sibling_type) 
		VALUES (?, ?, ?, ?, ?)
		`

	statement, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(
		rootHashID,
		self,
		parent,
		sibling,
		siblingType.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	logger.Logger.Debug(
		"added path to the database",
		zap.String("root hash", rootHash),
		zap.String("self", self),
		zap.String("parent", parent),
		zap.String("sibling", sibling),
		zap.String("siblingType", siblingType.String()),
	)

	return nil
}
