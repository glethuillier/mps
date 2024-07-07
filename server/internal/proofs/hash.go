package proofs

import (
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/glethuillier/fvs/server/internal/common"
	"github.com/glethuillier/fvs/server/internal/logger"
	"go.uber.org/zap"
)

// emptyHash returns the hash of an empty value ([]byte{})
func emptyHash(hashAlgorithm hash.Hash) string {
	hashAlgorithm.Reset()
	_, err := hashAlgorithm.Write([]byte{})
	if err != nil {
		logger.Logger.Panic(
			"cannot hash an empty value",
			zap.Error(err),
		)
	}
	return hex.EncodeToString(hashAlgorithm.Sum(nil))
}

func hashLeaf(hashAlgorithm hash.Hash, file *common.File) (string, error) {
	hashAlgorithm.Reset()
	_, err := hashAlgorithm.Write(file.Contents)
	if err != nil {
		return "", fmt.Errorf("cannot hash leaf: %w", err)
	}
	return hex.EncodeToString(hashAlgorithm.Sum(nil)), nil
}

func hashConcat(hashAlgorithm hash.Hash, a, b string) (string, error) {
	var err error
	hashAlgorithm.Reset()

	var hs [][]byte
	for _, h := range []string{a, b} {
		hBytes, err := hex.DecodeString(h)
		if err != nil {
			return "", err
		}

		hs = append(hs, hBytes)
	}

	for _, nodeHash := range hs {
		if _, err = hashAlgorithm.Write(nodeHash); err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(hashAlgorithm.Sum(nil)), nil
}
