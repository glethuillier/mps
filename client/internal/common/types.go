package common

import (
	"errors"

	"github.com/glethuillier/mps/lib/pkg/proofs"
)

type File struct {
	Filename string
	Contents []byte
	Proof    []proofs.ProofPart
	Error    error
}

type UploadRequest struct {
	Files []File
}

type DownloadRequest struct {
	ReceiptId string
	Filename  string
}

var ErrMismatchingRoots error = errors.New(
	"the request file is corrupted (root hashes do not match)",
)
