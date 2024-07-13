package client

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/glethuillier/mps/client/internal/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type connectionStatus = uint

const (
	NOT_CONNECTED connectionStatus = iota
	CONNECTING
	CONNECTED
)

type client struct {
	mu                 sync.RWMutex
	status             connectionStatus
	conn               *websocket.Conn
	url                url.URL
	messagesToSendC    chan interface{}
	messagesReceivedMC map[uuid.UUID]chan interface{}
}

// tryConnect attempts to connect to the server
// using a backoff strategy
func (c *client) tryConnect() {
	// if the client is already try to connect to the server,
	// just wait until it is connected
	if c.status == CONNECTING {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			<-ticker.C
			if c.conn == nil {
				continue
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				continue
			}
			return
		}
	}

	operation := func() error {
		logger.Logger.Debug("trying to connect to server")

		err := c.connect()
		if err != nil {
			c.status = CONNECTING
			logger.Logger.Debug(
				"reconnection attempt failed",
			)
			return err
		}

		c.status = CONNECTED
		return nil
	}

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = 1 * time.Second
	backoffConfig.MaxInterval = 10 * time.Second
	backoffConfig.MaxElapsedTime = 300 * time.Second

	err := backoff.Retry(operation, backoffConfig)
	if err != nil {
		logger.Logger.Fatal(
			"could not reconnect after maximum attempts",
			zap.Error(err),
		)
	}

	logger.Logger.Info("successfully connected")
}

func (c *client) connect() error {
	var err error

	c.conn, _, err = websocket.DefaultDialer.Dial(c.url.String(), nil)
	if err != nil {
		return err
	}

	return nil
}

// handle messages to be sent to the server
func (c *client) handleWrites(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		// regularly ping the server
		case <-ticker.C:
			if c.conn == nil {
				c.tryConnect()
				continue
			}
			err := c.send(websocket.PingMessage, []byte{})
			if err != nil {
				logger.Logger.Error(
					"write error",
					zap.Error(err),
				)
				c.tryConnect()
				continue
			}

		// send request to server
		case message := <-c.messagesToSendC:
			err := c.send(websocket.BinaryMessage, message.([]byte))
			if err != nil {
				logger.Logger.Error(
					"write close",
					zap.Error(err),
				)
				c.tryConnect()
				continue
			}

		case <-ctx.Done():
			err := c.send(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(
					websocket.CloseNormalClosure,
					"",
				))
			if err != nil {
				logger.Logger.Error(
					"write close",
					zap.Error(err),
				)
			}
			return
		}
	}
}

func (c *client) send(messageType int, message []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteMessage(messageType, message)
}

// handle messages received from the server
func (c *client) handleReads(ctx context.Context) {
	c.conn.SetPongHandler(func(appData string) error {
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			logger.Logger.Error(
				"error occurred while reading message from server",
				zap.Error(err),
			)
			c.tryConnect()
			continue
		}
		id, msg, err := processIncomingMessage(message)
		if err != nil {
			logger.Logger.Error(
				"error occurred while parsing message from server",
				zap.Error(err),
			)
		}

		c.mu.RLock()
		_, ok := c.messagesReceivedMC[id]
		if !ok {
			panic("chan map entry has not been initialized")
		}

		c.messagesReceivedMC[id] <- msg
		c.mu.RUnlock()
	}
}

func Run(
	ctx context.Context,
	messagesToSendC chan interface{},
	messagesReceivedMC map[uuid.UUID]chan interface{},
) {
	serverHost := os.Getenv("SERVER_HOST")
	if len(serverHost) == 0 {
		serverHost = "localhost"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if len(serverPort) == 0 {
		serverPort = "3000"
	}

	serverUrl := fmt.Sprintf("%s:%s", serverHost, serverPort)

	logger.Logger.Sugar().Infof("connecting to %s", serverUrl)

	client := &client{
		url: url.URL{
			Scheme: "ws",
			Host:   serverUrl,
			Path:   "/",
		},
		messagesToSendC:    messagesToSendC,
		messagesReceivedMC: messagesReceivedMC,
		status:             NOT_CONNECTED,
	}

	client.tryConnect()

	go client.handleReads(ctx)
	go client.handleWrites(ctx)
}
