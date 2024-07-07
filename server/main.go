package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/glethuillier/fvs/server/internal/logger"
	"github.com/glethuillier/fvs/server/internal/middleware"
	"github.com/glethuillier/fvs/server/internal/server"
	"go.uber.org/zap"
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
		logger.Logger.Info("closing the server...")
		cancel()
	}()

	// server <-> client
	requestsC := make(chan interface{}, 100)
	responsesC := make(chan interface{}, 100)
	go server.Run(requestsC, responsesC)

	err := middleware.Run(ctx, requestsC, responsesC)
	if err != nil {
		logger.Logger.Fatal(
			"middleware cannot be launched",
			zap.Error(err),
		)
	}
}
