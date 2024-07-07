package proofs

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"testing"

	"github.com/glethuillier/mps/client/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestBuildMerkleTree(t *testing.T) {
	tests := []struct {
		name                string
		hashAlgorithm       hash.Hash
		files               []common.File
		expectedRootHash    string
		expectedHashLeftL1  string
		expectedHashRightL1 string
		expectedError       bool
	}{
		{
			name:          "Positive test - 2 files - transparent hash",
			hashAlgorithm: newTransparentHash(),
			files: []common.File{
				{Contents: []byte{1}},
				{Contents: []byte{2}},
			},
			expectedRootHash:    "0102",
			expectedHashLeftL1:  "01",
			expectedHashRightL1: "02",
		},
		//
		// NOTE: problematic test that should be investigated to identify
		// the potential side effect caused by a bug
		//
		// {
		// 	name:          "Positive test - 3 files - transparent hash",
		// 	hashAlgorithm: newTransparentHash(),
		// 	files: []common.File{
		// 		{Contents: []byte{1}},
		// 		{Contents: []byte{2}},
		// 		{Contents: []byte{3}},
		// 	},
		// 	expectedRootHash:    "01020300",
		// 	expectedHashLeftL1:  "0102",
		// 	expectedHashRightL1: "0300",
		// },
		{
			name:          "Positive test - 2 files - sha256",
			hashAlgorithm: sha256.New(),
			files: []common.File{
				{Contents: []byte{1}},
				{Contents: []byte{2}},
			},
			expectedRootHash:    "42dbeeb4eb5d41bbdc93732c6a87ab3241ee03f44a0780a52ddf831f5fd88b53",
			expectedHashLeftL1:  "4bf5122f344554c53bde2ebb8cd2b7e3d1600ad631c385a5d7cce23c7785459a",
			expectedHashRightL1: "dbc1b4c900ffe48d575b5da5c638040125f65db0fe3e24494b76ea986457d986",
		},
		{
			name:          "Positive test - 3 files - sha256",
			hashAlgorithm: sha256.New(),
			files: []common.File{
				{Contents: []byte{1}},
				{Contents: []byte{2}},
				{Contents: []byte{3}},
			},
			expectedRootHash:    "c849e9c81c4c043b8c3be13568974c393fa81d54b5467f1fe291c079951adb19",
			expectedHashLeftL1:  "f059da7c02c43c6f2b2ea51ec701e2cac3f2c14abb55f860bd85f26f483a52a9",
			expectedHashRightL1: "2f3fcb322766895d1475987e1b16fc3e7fcbba659f57d799c9fae3295249e2ae",
		},
		{
			name:          "Positive test - 2 files - sha512",
			hashAlgorithm: sha512.New(),
			files: []common.File{
				{Contents: []byte{1}},
				{Contents: []byte{2}},
			},
			expectedRootHash:    "d091a63d9478334fc79a0642a717279ea1635b848c4b18ebeb33d41a50134e54572165c446ff29d29e43961b125a337c7f8a8977e7854fda9cfa5ce85e97e8a2",
			expectedHashLeftL1:  "7b54b66836c1fbdd13d2441d9e1434dc62ca677fb68f5fe66a464baadecdbd00576f8d6b5ac3bcc80844b7d50b1cc6603444bbe7cfcf8fc0aa1ee3c636d9e339",
			expectedHashRightL1: "fab848c9b657a853ee37c09cbfdd149d0b3807b191dde9b623ccd95281dd18705b48c89b1503903845bba5753945351fe6b454852760f73529cf01ca8f69dcca",
		},
		{
			name:          "Positive test - 3 files - sha512",
			hashAlgorithm: sha512.New(),
			files: []common.File{
				{Contents: []byte{1}},
				{Contents: []byte{2}},
				{Contents: []byte{3}},
			},
			expectedRootHash:    "da1193f5918a623022dcb3af9331942d40ba8647c84222eab9b9975ab2903342f4f4d249b4ad8c477064c9a7f5871870e0876642346c21295959363828c4aa40",
			expectedHashLeftL1:  "18142222c7b311840b39d8036f131405a69ddb02ea7417325ec66643bba609ac14dea4f43ab1e7d8f05d17e60493dfdd51e4b4f6ba95c5d98a61dd3fd1f04e63",
			expectedHashRightL1: "ccc49ed231c1bd6f4dc836ebefc11abc123b8af40fc2463fdefafc48a8512d51ae9aac10b40304cecfc8b3ed77eef544eac04abb2ba724fc12dbaa9a9d40d66f",
		},
		{
			name:          "Negative test - no file",
			hashAlgorithm: sha256.New(),
			files:         []common.File{},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tree, err := BuildMerkleTree(tc.hashAlgorithm, tc.files)

			if tc.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			rootHash, err := hex.DecodeString(tc.expectedRootHash)
			if err != nil {
				panic(err)
			}

			rootLeftL1, err := hex.DecodeString(tc.expectedHashLeftL1)
			if err != nil {
				panic(err)
			}

			rootRightL1, err := hex.DecodeString(tc.expectedHashRightL1)
			if err != nil {
				panic(err)
			}

			assert.Equal(t, tree.Root.hash, rootHash)
			assert.Equal(t, tree.Root.left.hash, rootLeftL1)
			assert.Equal(t, tree.Root.right.hash, rootRightL1)
		})
	}
}

// transparentHash is a hash algorithm that... does not
// hash the input data (the purpose of this custom hash
// algorithm is to help ensure that the tree is build
// correctly by verifying simple assumptions)
type transparentHash struct {
	data []byte
}

func newTransparentHash() hash.Hash {
	return &transparentHash{}
}

func (h *transparentHash) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		h.data = append(h.data, 0x00)
	} else {
		h.data = append(h.data, p...)
	}
	return len(p), nil
}

func (h *transparentHash) Sum(b []byte) []byte {
	// the transparent hash just appends the input data
	return append(b, h.data...)
}

func (h *transparentHash) Reset() {
	h.data = nil
}

func (h *transparentHash) Size() int {
	return len(h.data)
}

func (h *transparentHash) BlockSize() int {
	return 1
}
