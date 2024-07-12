package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/glethuillier/fvs/lib/pkg/proofs"
	"github.com/glethuillier/fvs/server/internal/common"
)

// IsTreeAlreadyPresent checks whether a root hash has already been
// saved in the database (meaning that the client already sent the
// files to the server). If a root hash is already present, the function
// returns the receipt ID corresponding to the root hash.
func (db *Database) IsTreeAlreadyPresent(rootHash string) (bool, string, error) {
	var receiptID string

	query := `SELECT receipt_id FROM RECEIPTS WHERE root_hash = ?`

	err := db.QueryRow(query, rootHash).Scan(&receiptID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, "", nil
		}
		return false, "", err
	}

	return true, receiptID, nil
}

// GetRootHash returns the root hash corresponding to a receipt ID
func (db *Database) GetRootHash(receiptId string) (string, error) {
	// get root hash
	var rootHash string
	query := "SELECT root_hash FROM RECEIPTS WHERE receipt_id = ?"
	err := db.QueryRow(query, receiptId).Scan(&rootHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("receipt_id %s not found", receiptId)
		}
		return "", err
	}

	return rootHash, nil
}

// GetTree returns a Merkle tree corresponding to a root hash
func (db *Database) GetTree(rootHash string) (*common.Tree, error) {
	tree := common.Tree{
		FilenameToHash: make(map[string]string),
		Nodes:          make(map[string]common.Node),
	}

	// get filenames and their hashes
	query := `
    SELECT
        f.filename,
        f.self_hash
    FROM
        FILES f
    JOIN
        RECEIPTS rh
    ON
        f.root_hash_id = rh.root_hash_id
    WHERE
        rh.root_hash = ?;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(rootHash)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var filename, selfHash string
		if err := rows.Scan(&filename, &selfHash); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		tree.FilenameToHash[filename] = selfHash
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	// get tree
	query = `
    SELECT
        p.self_hash,
        p.parent_hash,
        p.sibling_hash,
        p.sibling_type
    FROM
        RECEIPTS rh
    LEFT JOIN
        FILES f ON rh.root_hash_id = f.root_hash_id
    LEFT JOIN
        TREES p ON rh.root_hash_id = p.root_hash_id
    WHERE
        rh.root_hash = ?;`

	stmt, err = db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}
	defer stmt.Close()

	rows, err = stmt.Query(rootHash)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pathSelfHash, parentHash, siblingHash, siblingType sql.NullString
		if err := rows.Scan(
			&pathSelfHash,
			&parentHash,
			&siblingHash,
			&siblingType,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		tree.Nodes[pathSelfHash.String] = common.Node{
			Parent:      parentHash.String,
			Sibling:     siblingHash.String,
			SiblingType: proofs.GetSiblingType(siblingType.String),
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return &tree, nil
}
