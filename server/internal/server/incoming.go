package server

import (
	"github.com/glethuillier/fvs/lib/pkg/messages"
	"github.com/glethuillier/fvs/server/internal/common"
	"github.com/glethuillier/fvs/server/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// processIncomingMessage Protobuf deserializes the messages coming
// from the client
func processIncomingMessage(
	msg []byte,
	requestsC chan interface{},
) error {
	var wrapperMsg messages.WrapperMessage
	err := proto.Unmarshal(msg, &wrapperMsg)
	if err != nil {
		return err
	}

	requestId, err := uuid.Parse(wrapperMsg.MessageId)
	if err != nil {
		return err
	}

	switch wrapperMsg.Type {
	// preflight
	case messages.MessageType_TRANSFER_PREFLIGHT:
		var preflight messages.TransferPreflight
		err = proto.Unmarshal(wrapperMsg.Payload, &preflight)
		if err != nil {
			return err
		}

		logger.Logger.Debug(
			"received preflight",
			zap.Any("filenames", preflight.Filenames),
		)

		requestsC <- common.TransferRequest{
			MessageId: requestId,
			RootHash:  wrapperMsg.RootHash,
			Filenames: preflight.Filenames,
		}

	// receive file
	case messages.MessageType_TRANSFER_FILE:
		var receivedFile messages.TransferFile
		err = proto.Unmarshal(wrapperMsg.Payload, &receivedFile)
		if err != nil {
			return err
		}

		logger.Logger.Debug(
			"received file",
			zap.String("filename", receivedFile.Filename),
		)

		requestsC <- &common.File{
			MessageId: requestId,
			RootHash:  wrapperMsg.RootHash,
			Filename:  receivedFile.Filename,
			Contents:  receivedFile.Contents,
		}

	// send file
	case messages.MessageType_DOWNLOAD_REQUEST:
		var request messages.DownloadRequest
		err = proto.Unmarshal(wrapperMsg.Payload, &request)
		if err != nil {
			return err
		}

		logger.Logger.Debug(
			"received download request",
			zap.String("filename", request.Filename),
			zap.String("root_hash", request.RootHash),
		)

		requestsC <- common.DownloadRequest{
			MessageId: requestId,
			RootHash:  request.RootHash,
			Filename:  request.Filename,
		}

	default:
		logger.Logger.Error(
			"message from client cannot be deserialized (invalid type)",
		)
	}

	return nil
}
