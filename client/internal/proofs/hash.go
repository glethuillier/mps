package proofs

import (
	"fmt"
	"hash"
	"sync"

	"github.com/glethuillier/mps/client/internal/common"
	"github.com/glethuillier/mps/client/internal/logger"
	"go.uber.org/zap"
)

// EmptyHash returns the hash corresponding to
// an empty value ([]byte{})
type EmptyHash struct {
	value []byte
}

var (
	instance *EmptyHash
	once     sync.Once
)

type Hasher struct {
	hash.Hash
}

func GetHasher(hashAlgorithm hash.Hash) *Hasher {
	return &Hasher{
		hashAlgorithm,
	}
}

// NOTE: the client processes hashes as bytes

func emptyHash(hashAlgorithm hash.Hash) []byte {
	once.Do(func() {
		hashAlgorithm.Reset()
		_, err := hashAlgorithm.Write([]byte{})
		if err != nil {
			logger.Logger.Panic(
				"cannot hash an empty value",
				zap.Error(err),
			)
		}

		instance = &EmptyHash{
			value: hashAlgorithm.Sum(nil),
		}
	})

	return instance.value
}

func (h *Hasher) hashLeaf(file *common.File) ([]byte, error) {
	h.Reset()

	_, err := h.Write(file.Contents)
	if err != nil {
		return nil, fmt.Errorf("cannot hash leaf: %w", err)
	}

	return h.Sum(nil), nil
}

func (h *Hasher) hashConcat(a, b []byte) ([]byte, error) {
	var err error
	h.Reset()

	for _, nodeHash := range [][]byte{a, b} {
		if _, err = h.Write(nodeHash); err != nil {
			return nil, err
		}
	}

	return h.Sum(nil), nil
}
