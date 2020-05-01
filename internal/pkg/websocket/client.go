package websocket

import (
	"encoding/json"
	"fmt"
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
	Hub *Hub
}

type message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func (c *Client) wsClientReader() {
	defer func() {
		//unregister clients here
		log.Info().Msg("wsClientReader Connection Closed")
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
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

	for {
		_, rawMsgData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Warn().Err(err).Msg("Client closed ws connection unexpectedly")
			}
			log.Debug().Err(err).Msg("Client closed ws connection err: ")
			break
		}

		Msg := message{}

		//The web client sends ping messages as normal text
		if string(rawMsgData) != "Ping" {
			fmt.Println("Message Received!!:", string(rawMsgData))
			c.Hub.Broadcast <- rawMsgData

			//if err := json.Unmarshal(rawMsgData, &Msg); err != nil {
			//	log.Warn(err)
			//}
			//fmt.Println("Unmarshalled:")
			//fmt.Println(Msg)
		}

		Msg.Data = "Echo: " + string(rawMsgData)
		Msg.Type = "WSS"
		marshalledData, _ := json.Marshal(Msg)
		fmt.Println(string(marshalledData))

		c.Send <- marshalledData
		// broad cast message here
	}
}

func (c *Client) wsClientWriter() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.conn.Close()
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Debug().Msg("The hub closed the channel.")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Warn().Err(err).Msg("Could not write to websocket client")
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Warn().Err(err).Msg("could not send ping")
				return
			}
		}
	}
}
