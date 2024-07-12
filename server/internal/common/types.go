package common

import (
	"github.com/glethuillier/fvs/lib/pkg/proofs"
	"github.com/google/uuid"
)

type File struct {
	MessageId uuid.UUID
	RootHash  string
	Filename  string
	Contents  []byte
	Proof     []proofs.ProofPart
	Error     error
}

type TransferRequest struct {
	MessageId uuid.UUID
	RootHash  string
	Filenames []string
}

type DownloadRequest struct {
	MessageId uuid.UUID
	RootHash  string
	Filename  string
}

type TransferAck struct {
	MessageId uuid.UUID
	ReceiptId string
	Error     error
}

type ErrorResponse struct {
	MessageId uuid.UUID
	Error     error
}

// Merkle Tree

type Tree struct {
	RootHash string

	// filename -> self hash
	FilenameToHash map[string]string
	Nodes          map[string]Node
}
type Node struct {
	Parent      string
	SiblingType proofs.SiblingType
	Sibling     string
}
