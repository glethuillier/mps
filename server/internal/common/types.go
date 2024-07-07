package common

import "github.com/glethuillier/fvs/lib/pkg/proofs"

// messages from and to the client
type File struct {
	RequestId string
	Filename  string
	Contents  []byte
	Proof     []proofs.ProofPart
	Error     error
}

type TransferRequest struct {
	RequestId string
	Filenames []string
}

type DownloadRequest struct {
	ReceiptId string
	Filename  string
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
