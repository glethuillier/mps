package proofs

import (
	"crypto/sha512"
	"hash"
	"testing"

	"github.com/glethuillier/mps/client/internal/common"
	"github.com/glethuillier/mps/lib/pkg/proofs"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	tests := []struct {
		name             string
		hashAlgorithm    hash.Hash
		file             *common.File
		expectedRootHash string
		proof            []proofs.ProofPart
		expectedError    bool
	}{
		{
			name:          "Positive test - leaf 1/2",
			hashAlgorithm: sha512.New(),
			file: &common.File{
				Contents: []byte{1},
			},
			expectedRootHash: "d091a63d9478334fc79a0642a717279ea1635b848c4b18ebeb33d41a50134e54572165c446ff29d29e43961b125a337c7f8a8977e7854fda9cfa5ce85e97e8a2",
			proof: []proofs.ProofPart{
				{
					SiblingHash: "fab848c9b657a853ee37c09cbfdd149d0b3807b191dde9b623ccd95281dd18705b48c89b1503903845bba5753945351fe6b454852760f73529cf01ca8f69dcca",
					SiblingType: proofs.RightSibling,
				},
			},
		},
		{
			name:          "Positive test - leaf 2/2",
			hashAlgorithm: sha512.New(),
			file: &common.File{
				Contents: []byte{2},
			},
			expectedRootHash: "d091a63d9478334fc79a0642a717279ea1635b848c4b18ebeb33d41a50134e54572165c446ff29d29e43961b125a337c7f8a8977e7854fda9cfa5ce85e97e8a2",
			proof: []proofs.ProofPart{
				{
					SiblingHash: "7b54b66836c1fbdd13d2441d9e1434dc62ca677fb68f5fe66a464baadecdbd00576f8d6b5ac3bcc80844b7d50b1cc6603444bbe7cfcf8fc0aa1ee3c636d9e339",
					SiblingType: proofs.LeftSibling,
				},
			},
		},
		{
			name:          "Negative test - erroneous sibling hash",
			hashAlgorithm: sha512.New(),
			file: &common.File{
				Contents: []byte{1},
			},
			expectedRootHash: "d091a63d9478334fc79a0642a717279ea1635b848c4b18ebeb33d41a50134e54572165c446ff29d29e43961b125a337c7f8a8977e7854fda9cfa5ce85e97e8a2",
			proof: []proofs.ProofPart{
				{
					// hash: aa... instead of fa...
					SiblingHash: "aab848c9b657a853ee37c09cbfdd149d0b3807b191dde9b623ccd95281dd18705b48c89b1503903845bba5753945351fe6b454852760f73529cf01ca8f69dcca",
					SiblingType: proofs.RightSibling,
				},
			},
			expectedError: true,
		},
		{
			name:          "Negative test - inverse sibling",
			hashAlgorithm: sha512.New(),
			file: &common.File{
				Contents: []byte{1},
			},
			expectedRootHash: "d091a63d9478334fc79a0642a717279ea1635b848c4b18ebeb33d41a50134e54572165c446ff29d29e43961b125a337c7f8a8977e7854fda9cfa5ce85e97e8a2",
			proof: []proofs.ProofPart{
				{
					SiblingHash: "fab848c9b657a853ee37c09cbfdd149d0b3807b191dde9b623ccd95281dd18705b48c89b1503903845bba5753945351fe6b454852760f73529cf01ca8f69dcca",
					SiblingType: proofs.LeftSibling, // reverse sibling: right -> left
				},
			},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifyFile(tc.hashAlgorithm, tc.file, tc.expectedRootHash, tc.proof)

			if tc.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
