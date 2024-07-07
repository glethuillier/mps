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
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	logger.Init(logLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalC
		logger.Logger.Info("closing the client...")
		cancel()
	}()

	// caller -> client
	requestsC := make(chan interface{})
	responsesC := make(chan interface{})
	go server.Run(requestsC, responsesC)

	// client <-> server
	messagesToSendC := make(chan interface{})
	receivedMessagesC := make(chan interface{})
	go client.Run(ctx, messagesToSendC, receivedMessagesC)

	middleware.Run(
		ctx,
		requestsC,
		responsesC,
		messagesToSendC,
		receivedMessagesC,
	)
}
