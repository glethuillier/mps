package server

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/glethuillier/fvs/server/internal/common"
	"github.com/glethuillier/fvs/server/internal/logger"
	"github.com/glethuillier/fvs/server/internal/middleware"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// HandleConnections handles the connections between the server and the client
func HandleConnections(ctx context.Context, service *middleware.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Logger.Fatal(
				"cannot upgrade the connection to WebSockets",
				zap.Error(err),
			)
		}
		defer conn.Close()

		client := &Client{conn: conn, send: make(chan []byte)}

		requestsC := make(chan interface{})
		responsesC := make(chan interface{})

		service.Run(ctx, requestsC, responsesC)

		go client.handleReads(requestsC)
		client.handleWrites(responsesC)
	}
}

func (c *Client) handleReads(requestsC chan interface{}) {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			logger.Logger.Error("error reading message from client", zap.Error(err))
			break
		}

		if msgType == websocket.PingMessage {
			continue
		}

		if err := processIncomingMessage(msg, requestsC); err != nil {
			logger.Logger.Error(
				"cannot process message from client",
				zap.Error(err),
			)
		}
	}
}

func (c *Client) handleWrites(responsesC chan interface{}) {
	for {
		response := <-responsesC
		msg, err := prepareOutgoingMessage(response)
		if err != nil {
			logger.Logger.Error(
				"cannot prepare the message to send",
				zap.Error(err),
			)
			responsesC <- &common.TransferAck{
				Error:     err,
			}
		}

		err = c.conn.WriteMessage(
			websocket.BinaryMessage,
			msg,
		)
		if err != nil {
			logger.Logger.Error(
				"cannot prepare the message to send",
				zap.Error(err),
			)
			responsesC <- &common.TransferAck{
				Error: err,
			}
		}
	}
}

func Run(ctx context.Context, service *middleware.Service) {
	http.HandleFunc("/", HandleConnections(ctx, service))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	logger.Logger.Sugar().Infof("Server started on :%s", port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil)
	if err != nil {
		logger.Logger.Error("ListenAndServe: ", zap.Error(err))
	}
}
