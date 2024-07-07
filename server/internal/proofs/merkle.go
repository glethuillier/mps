package proofs

import (
	"crypto/sha512"
	"math"
	"sort"

	"github.com/glethuillier/fvs/lib/pkg/proofs"
	"github.com/glethuillier/fvs/server/internal/common"
)

// BuildMerkleTree builds the Merkle tree based on a list of files
// NOTE: iterative approach
func BuildMerkleTree(files []*common.File) (*common.Tree, error) {
	tree := common.Tree{
		FilenameToHash: make(map[string]string),
		Nodes:          make(map[string]common.Node),
	}

	// NOTE: should be configurable
	hashAlgorithm := sha512.New()

	// leaves
	leaves := []string{}
	for _, f := range files {
		h, err := hashLeaf(hashAlgorithm, f)
		if err != nil {
			return nil, err
		}
		tree.FilenameToHash[f.Filename] = h
		leaves = append(leaves, h)
	}

	// leaves nodes must be sorted to make the tree deterministic
	sort.Slice(leaves, func(i, j int) bool {
		return leaves[i] < leaves[j]
	})

	// n = 3; a = 4
	//
	//  a := power of 2 < n
	//  b := power of 2 > n
	//
	// 32 > 17 > 16
	//
	// (n / 2) - 1
	//
	// a : 16
	// b : 32
	//
	// b - (a + 1)

	// ensure that the number of leaves is a power of 2
	for {
		log2 := math.Log2(float64(len(leaves)))
		if log2 == float64(int(log2)) {
			break
		}

		leaves = append(leaves, emptyHash(hashAlgorithm))
	}

	level := leaves

	// compute root hash
	for {
		nextLevel := []string{}
		for i := 0; i < len(level); i += 2 {
			// NOTE: this approach is valid only because the tree
			// will be a perfect binary tree (as a consequence, the
			// number of nodes is necessary even for each level,
			// except for the last one, naturally).

			left := level[i]
			right := level[i+1]

			h, err := hashConcat(hashAlgorithm, left, right)
			if err != nil {
				return nil, err
			}

			tree.Nodes[right] = common.Node{
				Parent:      h,
				SiblingType: proofs.LeftSibling,
				Sibling:     left,
			}

			tree.Nodes[left] = common.Node{
				Parent:      h,
				SiblingType: proofs.RightSibling,
				Sibling:     right,
			}

			nextLevel = append(nextLevel, tree.Nodes[left].Parent)
		}

		level = nextLevel
		if len(level) == 1 {
			tree.Nodes[level[0]] = common.Node{}
			break
		}
	}

	tree.RootHash = level[0]

	return &tree, nil
}
