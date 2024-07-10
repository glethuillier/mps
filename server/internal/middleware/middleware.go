package middleware

import (
	"context"
	"fmt"

	"github.com/glethuillier/fvs/server/internal/common"
	"github.com/glethuillier/fvs/server/internal/database"
	"github.com/glethuillier/fvs/server/internal/helpers"
	"github.com/glethuillier/fvs/server/internal/logger"
	"github.com/glethuillier/fvs/server/internal/proofs"
	"go.uber.org/zap"
)

func Run(ctx context.Context, requestsC, responsesC chan interface{}) error {
	db, err := database.CreateDatabase("proofs.db")
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	err = helpers.Init()
	if err != nil {
		logger.Logger.Fatal("server cannot be initialized",
			zap.Error(err),
		)
	}

	receiver := receiver{
		db:            db,
		expectedFiles: make(map[string][]string),
	}

	filesC := make(chan *common.File)

	go receiver.receiveFiles(filesC, responsesC)

	for {
		select {
		case request := <-requestsC:
			switch r := request.(type) {

			case common.TransferRequest:
				receiver.prepareToReceiveFiles(r.RootHash, r.Filenames)

			case *common.File:
				filesC <- r

			case common.DownloadRequest:
				rootHash, err := db.GetRootHash(r.RootHash)
				if err != nil {
					logger.Logger.Error(
						"root hash cannot be retrieved from the database",
						zap.String("receipt_id", r.RootHash),
						zap.Error(err),
					)

					responsesC <- common.ErrorResponse{
						Error: err,
					}

					break
				}

				// get the requested file
				fileContents, err := helpers.GetFile(rootHash, r.Filename)
				if err != nil {
					logger.Logger.Error(
						"file contents cannot be retrieved",
						zap.String("receipt_id", r.Filename),
						zap.Error(err),
					)

					responsesC <- common.ErrorResponse{
						Error: fmt.Errorf("file not found"),
					}

					break
				}

				// get the tree
				tree, err := db.GetTree(rootHash)
				if err != nil {
					logger.Logger.Error(
						"the Merkle tree cannot be loaded",
						zap.String("root_hash", rootHash),
						zap.Error(err),
					)

					responsesC <- common.ErrorResponse{
						Error: err,
					}

					break
				}

				// extract the relevant proof from the tree
				proof, err := proofs.GenerateTransferableProof(tree, r.Filename)
				if err != nil {
					logger.Logger.Error(
						"proof cannot be communicated to the client",
						zap.Error(err),
					)

					responsesC <- common.ErrorResponse{
						Error: err,
					}

					break
				}

				responsesC <- &common.File{
					Filename: r.Filename,
					Contents: fileContents,
					Proof:    proof,
				}
			}

		case <-ctx.Done():
			return nil
		}
	}
}
