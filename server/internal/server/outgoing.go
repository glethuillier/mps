package server

import (
	"github.com/glethuillier/fvs/lib/pkg/messages"
	"github.com/glethuillier/fvs/lib/pkg/proofs"
	"github.com/glethuillier/fvs/server/internal/common"
	"github.com/glethuillier/fvs/server/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// prepareOutgoingMessage Protobuf serializes the messages to be
// sent to the client
func prepareOutgoingMessage(response interface{}) ([]byte, error) {
	var ack []byte
	var err error

	switch r := response.(type) {

	case common.TransferAck:
		if r.Error == nil {
			ack, err = proto.Marshal(&messages.TransferAck{
				StringOrArray: &messages.TransferAck_ReceiptId{
					ReceiptId: r.ReceiptId,
				},
			})
			if err != nil {
				return nil, err
			}
		} else {
			ack, err = proto.Marshal(&messages.TransferAck{
				StringOrArray: &messages.TransferAck_Error{
					Error: r.Error.Error(),
				},
			})
			if err != nil {
				return nil, err
			}
		}

		data, err := proto.Marshal(&messages.WrapperMessage{
			Type:    messages.MessageType_TRANSFER_ACK,
			Payload: ack,
		})
		if err != nil {
			logger.Logger.Error(
				"cannot marshal acknowledgment message",
				zap.Error(err),
			)
		}

		return data, nil

	// send file
	case *common.File:
		response, err := proto.Marshal(&messages.TransferFile{
			Filename: r.Filename,
			Contents: r.Contents,
			Proof:    encodeProof(r.Proof),
		})
		if err != nil {
			return nil, err
		}

		data, err := proto.Marshal(&messages.WrapperMessage{
			Type:    messages.MessageType_TRANSFER_FILE,
			Payload: response,
		})
		if err != nil {
			logger.Logger.Error(
				"cannot marshal file message",
				zap.Error(err),
			)
		}

		return data, nil

	// error
	case common.ErrorResponse:
		serverErr := r.Error.Error()
		response, err := proto.Marshal(&messages.TransferFile{
			Error: &serverErr,
		})
		if err != nil {
			return nil, err
		}

		data, err := proto.Marshal(&messages.WrapperMessage{
			Type:    messages.MessageType_TRANSFER_FILE,
			Payload: response,
		})
		if err != nil {
			logger.Logger.Error(
				"cannot marshal error message",
				zap.Error(err),
			)
		}

		return data, nil
	}

	return nil, nil
}

// encodeProof Protobuf serializes the proof
func encodeProof(proofParts []proofs.ProofPart) []*messages.ProofPart {
	var proof []*messages.ProofPart
	for _, p := range proofParts {
		proof = append(proof, &messages.ProofPart{
			SiblingType: messages.SiblingType(p.SiblingType),
			SiblingHash: p.SiblingHash,
		})
	}

	return proof
}
