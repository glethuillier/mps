package client

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/glethuillier/mps/client/internal/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type client struct {
	conn              *websocket.Conn
	url               url.URL
	messagesToSendC   chan interface{}
	receivedMessagesC chan interface{}
}

// tryConnect attempts to connect to the server
// using a backoff strategy
func (c *client) tryConnect(ctx context.Context) {
	operation := func() error {
		logger.Logger.Debug("trying to connect to server")

		err := c.connect()
		if err != nil {
			logger.Logger.Debug(
				"reconnection attempt failed",
			)
			return err
		}
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

	go c.handleReads(ctx)
	go c.handleWrites(ctx)

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
				return
			}
			err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				logger.Logger.Error(
					"write error",
					zap.Error(err),
				)
				c.tryConnect(ctx)
				return
			}

		// send request to server
		case message := <-c.messagesToSendC:
			err := c.conn.WriteMessage(websocket.BinaryMessage, message.([]byte))
			if err != nil {
				logger.Logger.Error(
					"write close",
					zap.Error(err),
				)
				c.tryConnect(ctx)
				return
			}

		case <-ctx.Done():
			err := c.conn.WriteMessage(
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

// handle messages received from the server
func (c *client) handleReads(ctx context.Context) {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			logger.Logger.Error(
				"error occurred while reading message from server",
				zap.Error(err),
			)
			c.tryConnect(ctx)
			return
		}
		msg, err := processIncomingMessage(message)
		if err != nil {
			logger.Logger.Error(
				"error occurred while parsing message from server",
				zap.Error(err),
			)
		}

		c.receivedMessagesC <- msg
	}
}

func Run(ctx context.Context, messagesToSendC, receivedMessagesC chan interface{}) {
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
		messagesToSendC:   messagesToSendC,
		receivedMessagesC: receivedMessagesC,
	}

	client.tryConnect(ctx)
}
