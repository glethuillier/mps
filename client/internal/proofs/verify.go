package proofs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/glethuillier/mps/client/internal/common"
	"github.com/glethuillier/mps/lib/pkg/proofs"
)

// VerifyFile verifies that a file has not been corrupted
// using a Merkle tree proof
func VerifyFile(
	hashAlgorithm hash.Hash,
	file *common.File,
	expectedRootHash string,
	proof []proofs.ProofPart,
) error {
	hasher := GetHasher(hashAlgorithm)

	// first, hash the file
	// (using the same hash algorithm used when it was uploaded)
	fileHash, err := hasher.hashLeaf(file)
	if err != nil {
		return err
	}

	current := fileHash

	// then verify the path to the root hash:
	// from the current node, reconstruct the parent using the sibling hash;
	// subsequently the parent becomes the current node, etc., up to the root
	for _, p := range proof {
		var (
			parent []byte
			err    error
		)

		siblingType, sibling := p.SiblingType, p.SiblingHash

		siblingHash, err := hex.DecodeString(sibling)
		if err != nil {
			return err
		}

		switch siblingType {

		case proofs.LeftSibling:
			parent, err = hasher.hashConcat(siblingHash, current)
			if err != nil {
				return err
			}

		case proofs.RightSibling:
			parent, err = hasher.hashConcat(current, siblingHash)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown sibling type: %s", siblingType)

		}

		current = parent
	}

	rootHashBytes, err := hex.DecodeString(expectedRootHash)
	if err != nil {
		return err
	}

	if !bytes.Equal(current, rootHashBytes) {
		return fmt.Errorf(
			"verification failed: (expected) %x != %x (actual)",
			rootHashBytes,
			current,
		)
	}

	return nil
}
