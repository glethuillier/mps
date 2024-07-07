package middleware

import (
	"context"
	"crypto/sha512"
	"fmt"
	"time"

	"github.com/glethuillier/mps/client/internal/client"
	"github.com/glethuillier/mps/client/internal/common"
	"github.com/glethuillier/mps/client/internal/database"
	"github.com/glethuillier/mps/client/internal/logger"
	"github.com/glethuillier/mps/client/internal/proofs"
	"go.uber.org/zap"
)

// WaitForConfirmation monitors incoming messages from
// the server and times out if the server has not sent
// an expected message in a given time
func WaitForConfirmation(
	ctx context.Context,
	receivedMessagesC chan interface{},
) (interface{}, error) {
	// NOTE: timeout should be configurable
	timer := time.After(60 * time.Second)

	for {
		select {
		case message := <-receivedMessagesC:
			return message, nil

		case <-timer:
			return nil,
				fmt.Errorf("the server has not processed all files")

		case <-ctx.Done():
			return nil, fmt.Errorf("internal server error")
		}
	}
}

func Run(
	ctx context.Context,
	requestsC,
	responsesC,
	messagesToSendC,
	receivedMessagesC chan interface{},
) error {
	sender := client.GetSender(messagesToSendC)

	db, err := database.CreateDatabase("roots.db")
	if err != nil {
		return fmt.Errorf("the database cannot be created: %w", err)
	}

	// NOTE: should be configurable
	hash := sha512.New()

	for {
		select {
		case request := <-requestsC:

			switch req := request.(type) {

			// upload files
			case common.UploadRequest:
				// build the Merkle tree
				tree, err := proofs.BuildMerkleTree(hash, req.Files)
				if err != nil {
					logger.Logger.Error("cannot build the tree",
						zap.Error(err))
				}

				id := tree.Root.GetHashAsString()

				sender.SendPreflightMessage(id, req)

				for _, f := range req.Files {
					go sender.SendFile(id, f)
				}

				// get the confirmation from the server
				response, err := WaitForConfirmation(ctx, receivedMessagesC)
				if err != nil {
					responsesC <- err
					break
				}

				switch resp := response.(type) {
				case error:
					responsesC <- resp
				case string:
					db.AddRootHash(resp, tree.Root.GetHash())
					responsesC <- resp
				}

				req.Files = nil

			// download files
			case common.DownloadRequest:
				sender.SendDownloadRequest(req.ReceiptId, req)

				// get the file from the server
				receivedFile, err := WaitForConfirmation(ctx, receivedMessagesC)
				if err != nil {
					responsesC <- err
					break
				}

				file, ok := receivedFile.(*common.File)
				if !ok || file.Error != nil {
					responsesC <- receivedFile
					break
				}

				// get the root hash corresponding to receipt ID
				rootHash, err := db.GetRootHash(req.ReceiptId)
				if err != nil {
					responsesC <- err
					break
				}

				// verify the proof
				verificationErr := proofs.VerifyFile(hash, file, rootHash, file.Proof)
				if verificationErr != nil {
					responsesC <- common.ErrMismatchingRoots
				} else {
					responsesC <- file
				}
			}

		case <-ctx.Done():
			return nil
		}
	}
}
