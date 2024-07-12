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
						10, 195, 169, 32, 112, 111, 117, 114, 32, 99, 101, 116, 32,
						101, 120, 101, 114, 99, 105, 99, 101, 32, 33, 32, 58, 41,
					},
				},
			},
			filesExpectedHashes: []string{
				"70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02",
				"eb5f59fd391278fa52091e4df383d12cfeaa815f5553a3afd567b250697c6ab5dcae6d8103cf1305a329c848e69ced500433716839ec7bf3af5b3f80a46bc486",
			},
			expectedRootHash: "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tree, err := BuildMerkleTree(tc.files)

			assert.NoError(t, err)
			assert.Equal(t, tree.RootHash, tc.expectedRootHash)

			for i, file := range tc.files {
				assert.Equal(t, tc.filesExpectedHashes[i], tree.FilenameToHash[file.Filename])
			}
		})
	}
}
