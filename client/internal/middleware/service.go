package middleware

import (
	"context"
	"crypto/sha512"
	"fmt"
	"hash"
	"sync"
	"time"

	"github.com/glethuillier/mps/client/internal/client"
	"github.com/glethuillier/mps/client/internal/common"
	"github.com/glethuillier/mps/client/internal/database"
	"github.com/glethuillier/mps/client/internal/logger"
	"github.com/glethuillier/mps/client/internal/proofs"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// receiveDataWithTimeout monitors incoming messages from
// the server and times out if the server has not sent
// an expected message in a given time
func receiveDataWithTimeout(
	ctx context.Context,
	messagesReceivedC chan interface{},
) (interface{}, error) {
	// TODO: timeout should be configurable
	timer := time.After(60 * time.Second)

	for {
		select {
		case message := <-messagesReceivedC:
			return message, nil

		case <-timer:
			return nil,
				fmt.Errorf("the server has not processed all files")

		case <-ctx.Done():
			return nil, fmt.Errorf("internal server error")
		}
	}
}

type Service struct {
	mu                sync.RWMutex
	db                *database.Database
	hashAlgorithm     hash.Hash
	sender            *client.Sender
	messagesReceivedC map[uuid.UUID]chan interface{}
}

func (s *Service) addReceiveChan(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messagesReceivedC[id] = make(chan interface{})
}

func (s *Service) deleteReceiveChan(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.messagesReceivedC, id)
}

func GetService(messagesToSendC chan interface{}, messagesReceivedC map[uuid.UUID]chan interface{}) (*Service, error) {
	sender := client.GetSender(messagesToSendC)

	db, err := database.CreateDatabase("roots.db")
	if err != nil {
		return nil, fmt.Errorf("the database cannot be created: %w", err)
	}

	// TODO: make it configurable
	hashAlgorithm := sha512.New()

	return &Service{
		db:                db,
		hashAlgorithm:     hashAlgorithm,
		sender:            sender,
		messagesReceivedC: messagesReceivedC,
	}, nil
}

func (s *Service) ProcessUploadRequest(
	ctx context.Context,
	requestId uuid.UUID,
	request common.UploadRequest,
) (string, error) {
	s.addReceiveChan(requestId)
	defer s.deleteReceiveChan(requestId)

	// build the Merkle tree
	tree, err := proofs.BuildMerkleTree(s.hashAlgorithm, request.Files)
	if err != nil {
		logger.Logger.Error("cannot build the tree",
			zap.Error(err))
	}

	rootHash := tree.Root.GetHashAsString()

	s.sender.SendPreflightMessage(requestId, rootHash, request)

	for _, f := range request.Files {
		go s.sender.SendFile(requestId, rootHash, f)
	}

	// get the confirmation from the server
	response, err := receiveDataWithTimeout(ctx, s.messagesReceivedC[requestId])
	if err != nil {
		return "", err
	}

	switch resp := response.(type) {
	case error:
		return "", resp
	case string:
		s.db.AddRootHash(resp, tree.Root.GetHash())
		return resp, nil
	}

	return "", nil
}

func (s *Service) ProcessDownloadRequest(
	ctx context.Context,
	requestId uuid.UUID,
	request common.DownloadRequest,
) (*common.File, error) {
	s.addReceiveChan(requestId)
	defer s.deleteReceiveChan(requestId)

	s.sender.SendDownloadRequest(requestId, request.ReceiptId, request)

	// get the file from the server
	data, err := receiveDataWithTimeout(ctx, s.messagesReceivedC[requestId])
	if err != nil {
		return nil, err
	}

	file, ok := data.(*common.File)
	if !ok {
		return nil, fmt.Errorf("data received from server is not a file: %T", file)
	}

	if file.Error != nil {
		return nil, file.Error
	}

	// get the root hash corresponding to receipt ID
	rootHash, err := s.db.GetRootHash(request.ReceiptId)
	if err != nil {
		return nil, err
	}

	// verify the proof
	verificationErr := proofs.VerifyFile(s.hashAlgorithm, file, rootHash, file.Proof)
	if verificationErr != nil {
		return nil, common.ErrMismatchingRoots
	} else {
		return file, nil
	}
}
