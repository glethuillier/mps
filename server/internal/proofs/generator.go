package proofs

import (
	"fmt"

	"github.com/glethuillier/fvs/lib/pkg/proofs"
	"github.com/glethuillier/fvs/server/internal/common"
)

// GenerateTransferableProof extract a proof from a tree corresponding to
// a given filename
func GenerateTransferableProof(tree *common.Tree, filename string) ([]proofs.ProofPart, error) {
	current, ok := tree.FilenameToHash[filename]
	if !ok {
		return nil, fmt.Errorf("filename %s not found in tree", filename)
	}

	var proofParts []proofs.ProofPart

	// starting from the file hash, identify the sibling at each level,
	// then go up to the root
	for {
		// if the current node has no sibling, that means it is the parent
		if tree.Nodes[current].SiblingType == proofs.NoSibling {
			return proofParts, nil
		}

		sibling := tree.Nodes[current].Sibling
		siblingType := tree.Nodes[current].SiblingType

		proofParts = append(proofParts, proofs.ProofPart{
			SiblingType: siblingType,
			SiblingHash: sibling,
		})

		current = tree.Nodes[current].Parent
	}
}
