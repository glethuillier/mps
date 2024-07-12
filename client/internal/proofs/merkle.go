package proofs

import (
	"encoding/hex"
	"fmt"
	"hash"
	"sort"

	"github.com/glethuillier/mps/client/internal/common"
)

type Tree struct {
	Hasher Hasher
	Root   *node
}

type node struct {
	hash  []byte
	left  *node
	right *node
}

// BuildMerkleTree builds a Merkle tree based from a list of files
func BuildMerkleTree(hashAlgorithm hash.Hash, files []common.File) (*Tree, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files to process")
	}

	t := Tree{
		Hasher: *GetHasher(hashAlgorithm),
	}

	nodes := make([]*node, 0)

	// create the nodes-leaves
	for _, f := range files {
		h, err := t.Hasher.hashLeaf(&f)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, &node{hash: h})
	}

	// leaves nodes must be sorted to make the tree deterministic
	sort.Slice(nodes, func(i, j int) bool {
		return string(nodes[i].hash) < string(nodes[j].hash)
	})

	// if the number of leaves is not a power of 2, fill in the leaves
	// needed to reach a power of 2 with empty data
	for {
		l := len(nodes)
		if l > 1 && (l&(l-1)) == 0 {
			break
		}

		nodes = append(nodes, &node{
			hash: emptyHash(hashAlgorithm),
		})
	}

	return t.BuildMerkleTree(nodes)
}

// BuildMerkleTree recursively builds the tree
func (t *Tree) BuildMerkleTree(nodes []*node) (*Tree, error) {
	// (base case) if one node remains, this is the root
	if len(nodes) == 1 {
		t.Root = nodes[0]
		return t, nil
	}

	var nextLevel []*node
	for i := 0; i < len(nodes); i += 2 {
		// NOTE: this approach is valid only because the tree
		// will be a perfect binary tree (as a consequence, the
		// number of nodes is necessary even for each level,
		// except for the last one, naturally).
		left := nodes[i]
		right := nodes[i+1]

		h, err := t.Hasher.hashConcat(left.hash, right.hash)
		if err != nil {
			return nil, err
		}

		nextLevel = append(nextLevel, &node{
			hash:  h,
			left:  left,
			right: right,
		})
	}

	return t.BuildMerkleTree(nextLevel)
}

func (n *node) GetHash() []byte {
	return n.hash
}

func (n *node) GetHashAsString() string {
	return hex.EncodeToString(n.hash)
}
