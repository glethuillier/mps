package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/glethuillier/mps/client/internal/client"
	"github.com/glethuillier/mps/client/internal/logger"
	"github.com/glethuillier/mps/client/internal/middleware"
	"github.com/glethuillier/mps/client/internal/server"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	logger.Init(logLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})

	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalC
		logger.Logger.Info("closing the client...")
		cancel()
		done <- struct{}{}
	}()

	messagesToSendC := make(chan interface{})

	// messages received from the mps server are
	// dispatched by request ID (map of channels)
	messagesReceivedMC := make(map[uuid.UUID]chan interface{})

	service, err := middleware.GetService(messagesToSendC, messagesReceivedMC)
	if err != nil {
		logger.Logger.Panic("cannot run the middleware", zap.Error(err))
	}

	// caller -> client
	go server.Run(ctx, service)

	// client <-> server
	go client.Run(ctx, messagesToSendC, messagesReceivedMC)

	<-done
}
