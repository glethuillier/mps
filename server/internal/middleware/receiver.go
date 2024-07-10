package middleware

import (
	"fmt"
	"sync"

	"github.com/glethuillier/fvs/server/internal/common"
	"github.com/glethuillier/fvs/server/internal/database"
	"github.com/glethuillier/fvs/server/internal/helpers"
	"github.com/glethuillier/fvs/server/internal/logger"
	"github.com/glethuillier/fvs/server/internal/proofs"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type responseType int

const (
	ROOTS_MATCH responseType = iota
	ROOTS_MISMATCH
	NOT_UNIQUE
	OTHER_ERROR
)

type receiver struct {
	sync.RWMutex
	db            *database.Database
	expectedFiles map[string][]string
}

func (r *receiver) prepareToReceiveFiles(requestId string, filenames []string) {
	r.Lock()
	defer r.Unlock()

	r.expectedFiles[requestId] = filenames
}

func (r *receiver) receiveFiles(
	filesC chan *common.File,
	responsesC chan interface{},
) bool {
	filesToProcess := make(map[string][]*common.File)

	for {
		select {
		case file := <-filesC:
			r.Lock()
			_, ok := filesToProcess[file.RootHash]
			if !ok {
				filesToProcess[file.RootHash] = []*common.File{}
			}

			filesToProcess[file.RootHash] = append(filesToProcess[file.RootHash], file)
			r.Unlock()

			r.RLock()
			if len(filesToProcess[file.RootHash]) == len(r.expectedFiles[file.RootHash]) {
				r.processFiles(file.RootHash, filesToProcess[file.RootHash], responsesC)

				// discard files in memory
				filesToProcess[file.RootHash] = nil
			}
			r.RUnlock()
		}
	}
}

func (r *receiver) processFiles(
	expectedRootHash string,
	files []*common.File,
	responsesC chan interface{},
) {
	var (
		responseType   responseType
		knownReceiptId string
	)

	tree, err := proofs.BuildMerkleTree(files)
	if err != nil {
		responseType = OTHER_ERROR

		logger.Logger.Error(
			"error occurred while computing the Merkle root hash",
			zap.Error(err),
		)
	} else {
		treeAlreadyPresent, receiptId, err := r.db.IsTreeAlreadyPresent(tree.RootHash)
		if treeAlreadyPresent {
			responseType = NOT_UNIQUE
			knownReceiptId = receiptId
			logger.Logger.Error(
				"tree already exists in database",
				zap.String("root_hash", tree.RootHash),
				zap.String("receipt_id", receiptId),
				zap.Error(err),
			)
		} else if tree.RootHash != expectedRootHash {
			responseType = ROOTS_MISMATCH
			logger.Logger.Error(
				"root hashes mismatch",
				zap.String("hash_from_client", expectedRootHash),
				zap.String("hash_from_server", tree.RootHash),
			)
		} else {
			responseType = ROOTS_MATCH
		}
	}

	switch responseType {
	case ROOTS_MATCH:
		for _, f := range files {
			helpers.SaveFile(f.RootHash, f.Filename, f.Contents)
		}

		receiptId := uuid.New()

		err = r.db.SaveTree(receiptId, tree)
		if err != nil {
			logger.Logger.Error(
				"the tree cannot be saved in database",
				zap.String("receipt_id", receiptId.String()),
				zap.Error(err),
			)

			responsesC <- common.TransferAck{
				Error: err,
			}
		} else {
			responsesC <- common.TransferAck{
				ReceiptId: receiptId.String(),
			}
		}

	case NOT_UNIQUE:
		responsesC <- common.TransferAck{
			Error: fmt.Errorf(
				"these files has already been processed by the server; receipt ID: %s",
				knownReceiptId,
			),
		}

	case OTHER_ERROR:
		responsesC <- common.TransferAck{
			// send a generic error message to the client so that we
			// do not leak detail about the internal implementation
			Error: fmt.Errorf(
				"the server cannot process the files",
			),
		}

	case ROOTS_MISMATCH:
		responsesC <- common.TransferAck{
			Error: fmt.Errorf(
				// send a generic error message to client (i.e., do not
				// specify that this information corresponds to Merkle
				// tree root hashes).
				//
				// NOTE: the client should in fact handle the conversion of
				// technical error messages into human-friendly error
				// messages
				"proofs mismatch: %s v. %s",
				expectedRootHash,
				tree.RootHash,
			),
		}
	}
}
