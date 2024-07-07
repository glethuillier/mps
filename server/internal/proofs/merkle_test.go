package proofs

import (
	"crypto/sha256"
	"hash"
	"testing"

	"github.com/glethuillier/fvs/server/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestBuildMerkleTree(t *testing.T) {
	tests := []struct {
		name                string
		hashAlgorithm       hash.Hash
		files               []*common.File
		filesExpectedHashes []string
		expectedRootHash    string
		expectedError       bool
	}{
		{
			name:          "Positive test - 2 files - sha256",
			hashAlgorithm: sha256.New(),
			files: []*common.File{
				{
					Filename: "readme.txt",
					Contents: []byte{
						89, 111, 117, 32, 97, 99, 116, 117, 97, 108, 108, 121, 32, 114,
						101, 97, 100, 32, 105, 116, 33,
					},
				},
				{
					Filename: "abc.txt",
					Contents: []byte{
						74, 39, 97, 105, 32, 116, 111, 117, 116, 32, 100, 111, 110,
						110, 195, 169, 32, 112, 111, 117, 114, 32, 99, 101, 116, 32,
						101, 120, 101, 114, 99, 105, 99, 101, 32, 33, 32, 58, 41,
					},
				},
			},
			filesExpectedHashes: []string{
				"70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02",
				"48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84",
			},
			expectedRootHash: "4dfefda1ce283a40b492f1f1e7d234f9e379ecdafaf172cf5a510ef4889bc9697bc63c2426c2456fc628e05db7277baee684a6a844cda178e7943084ea0dc29e",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tree, err := BuildMerkleTree(tc.files)

			assert.NoError(t, err)
			assert.Equal(t, tree.RootHash, tc.expectedRootHash)

			for i, file := range tc.files {
				assert.Equal(t, tree.FilenameToHash[file.Filename], tc.filesExpectedHashes[i])
			}
		})
	}
}
