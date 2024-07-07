package client

import (
	"github.com/glethuillier/mps/client/internal/common"
	"github.com/glethuillier/mps/client/internal/logger"
	"go.uber.org/zap"

	"github.com/glethuillier/mps/lib/pkg/messages"
	"google.golang.org/protobuf/proto"
)

type Sender struct {
	messagesC chan interface{}
}

// GetSender returns a Sender which handles messages to
// be sent to the server
func GetSender(messagesC chan interface{}) *Sender {
	return &Sender{
		messagesC: messagesC,
	}
}

// SendPreflightMessage Protobuf serializes preflight messages
func (s *Sender) SendPreflightMessage(id string, request common.UploadRequest) {
	var filenames []string
	for _, f := range request.Files {
		filenames = append(filenames, f.Filename)
	}

	init, err := proto.Marshal(&messages.TransferPreflight{
		Filenames:     filenames,
		HashAlgorithm: messages.HashAlgorithm_SHA512,
	})
	if err != nil {
		logger.Logger.Error(
			"cannot marshal preflight message",
			zap.Error(err),
		)
	}

	data, err := proto.Marshal(&messages.WrapperMessage{
		RequestId: id,
		Type:      messages.MessageType_TRANSFER_PREFLIGHT,
		Payload:   init,
	})
	if err != nil {
		logger.Logger.Error(
			"cannot marshal preflight message in wrapper",
			zap.Error(err),
		)
	}

	s.messagesC <- data
}

// SendDownloadRequest Protobuf serializes requests to download files
func (s *Sender) SendDownloadRequest(id string, request common.DownloadRequest) {
	req, err := proto.Marshal(&messages.DownloadRequest{
		ReceiptId: request.ReceiptId,
		Filename:  request.Filename,
	})
	if err != nil {
		logger.Logger.Error(
			"cannot marshal download request",
			zap.Error(err),
		)
	}

	data, err := proto.Marshal(&messages.WrapperMessage{
		RequestId: id,
		Type:      messages.MessageType_DOWNLOAD_REQUEST,
		Payload:   req,
	})
	if err != nil {
		logger.Logger.Error(
			"cannot marshal download request in wrapper",
			zap.Error(err),
		)
	}

	s.messagesC <- data
}

// SendFile Protobuf serializes files to be sent to the server
func (s *Sender) SendFile(id string, request common.File) {
	init, err := proto.Marshal(&messages.TransferFile{
		Filename: request.Filename,
		Contents: request.Contents,
	})
	if err != nil {
		logger.Logger.Error(
			"cannot marshal file",
			zap.Error(err),
		)
	}

	data, err := proto.Marshal(&messages.WrapperMessage{
		RequestId: id,
		Type:      messages.MessageType_TRANSFER_FILE,
		Payload:   init,
	})
	if err != nil {
		logger.Logger.Error(
			"cannot marshal download request in wrapper",
			zap.Error(err),
		)
	}

	s.messagesC <- data
}
