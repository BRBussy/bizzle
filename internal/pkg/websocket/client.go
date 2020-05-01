package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 4096
)

type Client struct {
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	Send chan []byte
	// hub
	// Hub *Hub
}

type message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func (c *Client) wsClientReader() {
	// close websocket on termination of client reader
	defer func() {
		log.Info().Msg("wsClientReader Connection Closed")
		if err := c.conn.Close(); err != nil {
			log.Error().Err(err).Msg("error closing websocket client connection")
		}
	}()

	// setup connection parameters
	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Error().Err(err).Msg("error setting maximum read deadline on websocket")
		return
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Error().Err(err).Msg("error setting maximum read deadline on websocket")
			return exception.ErrUnexpected{}
		}
		return nil
	})

	// Send introductory message to client
	helloClient, err := json.Marshal(message{
		Type: "WSS",
		Data: "hello client!",
	})
	if err == nil {
		c.Send <- helloClient
	} else {
		log.Warn().Msg("Unable to marshal message for client")
	}

	// go into read monitoring
	for {
		_, rawMsgData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Warn().Err(err).Msg("Client closed ws connection unexpectedly")
			}
			log.Debug().Err(err).Msg("Client closed ws connection err: ")
			break
		}
		fmt.Println("Message Received!!:", string(rawMsgData))
	}
}

func (c *Client) wsClientWriter() {
	// create ticker to ping socket connection
	pingTicker := time.NewTicker(pingPeriod)
	fmt.Println(pingPeriod)

	// close connection on return of client writer
	defer func() {
		if err := c.conn.Close(); err != nil {
			log.Error().Err(err).Msg("error closing websocket client connection")
		}
		pingTicker.Stop()
	}()

	// go into write monitor loop
	for {
		select {

		// monitor send client channel for outgoing messages to transmit
		case message, ok := <-c.Send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Error().Err(err).Msg("error setting write deadline for send message")
				return
			}
			if !ok {
				log.Debug().Msg("send channel closed")
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Error().Err(err).Msg("error sending close message")
				}
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Warn().Err(err).Msg("could not write to websocket")
				return
			}

		// keep websocket connection alive by sending ping messages
		case <-pingTicker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Error().Err(err).Msg("error setting write deadline for ping message")
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Warn().Err(err).Msg("could not send ping")
				return
			}
		}
	}
}
