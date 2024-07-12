package proofs

import (
	"reflect"
	"testing"

	"github.com/glethuillier/fvs/lib/pkg/proofs"
	"github.com/glethuillier/fvs/server/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	tests := []struct {
		name             string
		files            []*common.File
		tree             *common.Tree
		expectedProof    [][]proofs.ProofPart
		expectedMismatch bool
	}{
		{
			name: "Positive test - 2 files",
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
						110, 195, 169, 32, 112, 111, 117, 114, 32, 99, 101, 116, 32,
						101, 120, 101, 114, 99, 105, 99, 101, 32, 33, 32, 58, 41,
					},
				},
			},
			tree: &common.Tree{
				RootHash: "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
				FilenameToHash: map[string]string{
					"abc.txt":    "48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84",
					"readme.txt": "70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02",
				},
				Nodes: map[string]common.Node{
					"040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649": {SiblingType: proofs.NoSibling},
					"48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84": {
						Parent:      "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
						Sibling:     "70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02",
						SiblingType: proofs.RightSibling,
					},
					"70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02": {
						Parent:      "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
						Sibling:     "48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84",
						SiblingType: proofs.LeftSibling,
					},
				},
			},
			expectedProof: [][]proofs.ProofPart{
				{
					{
						SiblingHash: "48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84",
						SiblingType: proofs.LeftSibling,
					},
				},
				{
					{
						SiblingHash: "70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02",
						SiblingType: proofs.RightSibling,
					},
				},
			},
		},
		{
			name: "Positive test - 3 files",
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
						110, 195, 169, 32, 112, 111, 117, 114, 32, 99, 101, 116, 32,
						101, 120, 101, 114, 99, 105, 99, 101, 32, 33, 32, 58, 41,
					},
				},
				{
					Filename: "img.png",
					Contents: []byte{
						1, 2, 3, 4, 5, 6, 7, 8, 9,
					},
				},
			},
			tree: &common.Tree{
				RootHash: "2c6d246cb5fc78916395cd786aa0264c5211e65bfe03724dc6a264b3e53b0960107b2b9629f7a17ab467313a20a5cc513558444c8a9013a879d567c273c5c748",
				FilenameToHash: map[string]string{
					"abc.txt":    "48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84",
					"readme.txt": "70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02",
					"img.png":    "c99f95c6d2b8ed5a065946e0e6a4f76a2e7bbd5fe4d8ca922b7f1537dbf14db24145403f736a689d15a5c5e72d2742bee85420e54f7813439125273f112d133a",
				},
				Nodes: map[string]common.Node{
					"2c6d246cb5fc78916395cd786aa0264c5211e65bfe03724dc6a264b3e53b0960107b2b9629f7a17ab467313a20a5cc513558444c8a9013a879d567c273c5c748": {SiblingType: proofs.NoSibling},
					"48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84": {
						Parent:      "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
						Sibling:     "70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02",
						SiblingType: proofs.RightSibling,
					},
					"70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02": {
						Parent:      "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
						Sibling:     "48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84",
						SiblingType: proofs.LeftSibling,
					},
					"040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649": {
						Parent:      "2c6d246cb5fc78916395cd786aa0264c5211e65bfe03724dc6a264b3e53b0960107b2b9629f7a17ab467313a20a5cc513558444c8a9013a879d567c273c5c748",
						Sibling:     "64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532",
						SiblingType: proofs.RightSibling,
					},
					"64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532": {
						Parent:      "2c6d246cb5fc78916395cd786aa0264c5211e65bfe03724dc6a264b3e53b0960107b2b9629f7a17ab467313a20a5cc513558444c8a9013a879d567c273c5c748",
						Sibling:     "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
						SiblingType: proofs.LeftSibling,
					},
					"c99f95c6d2b8ed5a065946e0e6a4f76a2e7bbd5fe4d8ca922b7f1537dbf14db24145403f736a689d15a5c5e72d2742bee85420e54f7813439125273f112d133a": {
						Parent:      "64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532",
						Sibling:     "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
						SiblingType: proofs.RightSibling,
					},
					"cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e": {
						Parent:      "64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532",
						Sibling:     "c99f95c6d2b8ed5a065946e0e6a4f76a2e7bbd5fe4d8ca922b7f1537dbf14db24145403f736a689d15a5c5e72d2742bee85420e54f7813439125273f112d133a",
						SiblingType: proofs.LeftSibling,
					},
				},
			},
			expectedProof: [][]proofs.ProofPart{
				{
					{
						SiblingHash: "48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84",
						SiblingType: proofs.LeftSibling,
					},
					{
						SiblingHash: "64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532",
						SiblingType: proofs.RightSibling,
					},
				},
				{
					{
						SiblingHash: "70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02",
						SiblingType: proofs.RightSibling,
					},
					{
						SiblingHash: "64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532",
						SiblingType: proofs.RightSibling,
					},
				},
				{
					{
						SiblingHash: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
						SiblingType: proofs.RightSibling,
					},
					{
						SiblingHash: "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
						SiblingType: proofs.LeftSibling,
					},
				},
			},
		},
		{
			name: "Negative test - 3 files - wrong sibling type",
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
						110, 195, 169, 32, 112, 111, 117, 114, 32, 99, 101, 116, 32,
						101, 120, 101, 114, 99, 105, 99, 101, 32, 33, 32, 58, 41,
					},
				},
				{
					Filename: "img.png",
					Contents: []byte{
						1, 2, 3, 4, 5, 6, 7, 8, 9,
					},
				},
			},
			tree: &common.Tree{
				RootHash: "2c6d246cb5fc78916395cd786aa0264c5211e65bfe03724dc6a264b3e53b0960107b2b9629f7a17ab467313a20a5cc513558444c8a9013a879d567c273c5c748",
				FilenameToHash: map[string]string{
					"abc.txt":    "48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84",
					"readme.txt": "70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02",
					"img.png":    "c99f95c6d2b8ed5a065946e0e6a4f76a2e7bbd5fe4d8ca922b7f1537dbf14db24145403f736a689d15a5c5e72d2742bee85420e54f7813439125273f112d133a",
				},
				Nodes: map[string]common.Node{
					"2c6d246cb5fc78916395cd786aa0264c5211e65bfe03724dc6a264b3e53b0960107b2b9629f7a17ab467313a20a5cc513558444c8a9013a879d567c273c5c748": {SiblingType: proofs.NoSibling},
					"48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84": {
						Parent:      "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
						Sibling:     "70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02",
						SiblingType: proofs.RightSibling,
					},
					"70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02": {
						Parent:      "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
						Sibling:     "48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84",
						SiblingType: proofs.LeftSibling,
					},
					"040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649": {
						Parent:      "2c6d246cb5fc78916395cd786aa0264c5211e65bfe03724dc6a264b3e53b0960107b2b9629f7a17ab467313a20a5cc513558444c8a9013a879d567c273c5c748",
						Sibling:     "64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532",
						SiblingType: proofs.RightSibling,
					},
					"64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532": {
						Parent:      "2c6d246cb5fc78916395cd786aa0264c5211e65bfe03724dc6a264b3e53b0960107b2b9629f7a17ab467313a20a5cc513558444c8a9013a879d567c273c5c748",
						Sibling:     "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
						SiblingType: proofs.LeftSibling,
					},
					"c99f95c6d2b8ed5a065946e0e6a4f76a2e7bbd5fe4d8ca922b7f1537dbf14db24145403f736a689d15a5c5e72d2742bee85420e54f7813439125273f112d133a": {
						Parent:      "64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532",
						Sibling:     "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
						SiblingType: proofs.RightSibling,
					},
					"cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e": {
						Parent:      "64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532",
						Sibling:     "c99f95c6d2b8ed5a065946e0e6a4f76a2e7bbd5fe4d8ca922b7f1537dbf14db24145403f736a689d15a5c5e72d2742bee85420e54f7813439125273f112d133a",
						SiblingType: proofs.LeftSibling,
					},
				},
			},
			expectedProof: [][]proofs.ProofPart{
				{
					{
						SiblingHash: "48e69af2e737b5e6ebd3c129838b8b582bd7bdbdb6ec1c6e99ed311031b3735a819ca2ef5c68bf054891f9ab1928bcae851e943b03d2cd0842ce40b4bc9ceb84",
						SiblingType: proofs.RightSibling, // here: left -> right (wrong sibling type)
					},
					{
						SiblingHash: "64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532",
						SiblingType: proofs.RightSibling,
					},
				},
				{
					{
						SiblingHash: "70a49087db423f89aeea154a0f961f4aef0e634b286e3fdf35b430403421f031daf301ec0da455e226bcba40720f2147cbb7fa638917ee67a8fc40b143fa5c02",
						SiblingType: proofs.RightSibling,
					},
					{
						SiblingHash: "64b7683abaa13b47aa21928b7a43dc6e7459701d49f111fffb1ee0abca339db32cdee40746b2834b52063e90587f5f645c0b4a58247d1bcefbaa52ab95af3532",
						SiblingType: proofs.LeftSibling, // here: right -> left (wrong sibling type)
					},
				},
				{
					{
						SiblingHash: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
						SiblingType: proofs.LeftSibling, // here: right -> left (wrong sibling type)
					},
					{
						SiblingHash: "040038907ccb5294981ecd6c653a7c1528844ccfb6a3c62d61f7485c9afc762d18ceefb9ebf873d8e2a3cc656796e0130a8546adced12952772deed871bef649",
						SiblingType: proofs.LeftSibling,
					},
				},
			},
			expectedMismatch: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for i, p := range tc.files {
				proof, err := GenerateTransferableProof(tc.tree, p.Filename)

				if tc.expectedMismatch {
					assert.False(t, reflect.DeepEqual(tc.expectedProof[i], proof))
					continue
				}

				assert.NoError(t, err)
				assert.True(t, reflect.DeepEqual(tc.expectedProof[i], proof))
			}
		})
	}
}
