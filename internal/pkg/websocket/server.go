package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout:  0,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	WriteBufferPool:   nil,
	Subprotocols:      nil,
	Error:             nil,
	CheckOrigin:       func(r *http.Request) bool { return true },
	EnableCompression: false,
}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, hub *Hub) {
	// allow any origin to connect

	// upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Could not upgrade ws connection")
		w.WriteHeader(500)
		return
	}

	// construct a new client
	newClient := &Client{
		conn: conn,
		Send: make(chan []byte),
		Hub:  hub,
	}

	// Register Client with hub
	hub.Register <- newClient

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines
	go newClient.wsClientReader()
	go newClient.wsClientWriter()
}
