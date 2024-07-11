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

	done := make(chan struct{})

	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalC
		logger.Logger.Info("closing the server...")
		cancel()
		done <- struct{}{}
	}()

	service, err := middleware.GetService()
	if err != nil {
		logger.Logger.Fatal(
			"middleware cannot be launched",
			zap.Error(err),
		)
	}

	go server.Run(ctx, service)

	<-done
}
