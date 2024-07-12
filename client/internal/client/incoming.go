package client

import (
	"fmt"

	"github.com/glethuillier/mps/client/internal/common"
	"github.com/glethuillier/mps/client/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/glethuillier/mps/lib/pkg/messages"
	"github.com/glethuillier/mps/lib/pkg/proofs"
	"google.golang.org/protobuf/proto"
)

// processIncomingMessage parses individual raw messages from the server
func processIncomingMessage(msg []byte) (uuid.UUID, interface{}, error) {
	var wrapperMsg messages.WrapperMessage

	err := proto.Unmarshal(msg, &wrapperMsg)
	if err != nil {
		return uuid.UUID{}, nil, err
	}

	id, err := uuid.Parse(wrapperMsg.MessageId)
	if err != nil {
		return uuid.UUID{}, nil, err
	}

	switch wrapperMsg.Type {

	// ack
	case messages.MessageType_TRANSFER_ACK:
		var ack messages.TransferAck
		err = proto.Unmarshal(wrapperMsg.Payload, &ack)
		if err != nil {
			return id, nil, err
		}

		logger.Logger.Debug(
			"received ack",
			zap.Any("message", &ack),
		)

		serverErr := ack.GetError()
		if serverErr != "" {
			return id, fmt.Errorf(serverErr), nil
		}

		return id, ack.GetReceiptId(), nil

	// receive file
	case messages.MessageType_TRANSFER_FILE:
		var file messages.TransferFile
		err = proto.Unmarshal(wrapperMsg.Payload, &file)
		if err != nil {
			return id, nil, err
		}

		if file.Error != nil {
			return id, &common.File{
				Error: fmt.Errorf("error returned by server: %s", *file.Error),
			}, nil
		} else {
			return id, &common.File{
				Filename: file.Filename,
				Contents: file.Contents,
				Proof:    deserializeProof(file.Proof),
			}, nil
		}
	}

	return uuid.UUID{}, nil, nil
}

// deserializeProof transforms a Protobuf serialized proof sent by the server
// into a proof part object
func deserializeProof(proofPaths []*messages.ProofPart) []proofs.ProofPart {
	var paths []proofs.ProofPart

	for _, p := range proofPaths {
		paths = append(paths, proofs.ProofPart{
			SiblingType: proofs.SiblingType(p.SiblingType),
			SiblingHash: p.SiblingHash,
		})
	}

	return paths
}
