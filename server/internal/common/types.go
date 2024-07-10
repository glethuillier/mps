package common

import "github.com/glethuillier/fvs/lib/pkg/proofs"

// messages from and to the client
type File struct {
	RootHash string
	Filename string
	Contents []byte
	Proof    []proofs.ProofPart
	Error    error
}

type TransferRequest struct {
	RootHash  string
	Filenames []string
}

type DownloadRequest struct {
	RootHash string
	Filename string
}

type TransferAck struct {
	ReceiptId string
	Error     error
}

type ErrorResponse struct {
	Error error
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
