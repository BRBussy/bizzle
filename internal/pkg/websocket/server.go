package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, hub *Hub) {
	fmt.Println("Client Connected!")
	//Allow any origin to connect
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

	//Upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Could not upgrade ws connection")
		w.WriteHeader(500)
		return
	}

	//Construct a new client
	newClient := &Client{
		conn: conn,
		Send: make(chan []byte),
		Hub:  hub,
	}

	fmt.Println("Register client!")
	// Register Client with hub
	hub.Register <- newClient
	fmt.Println("Done Register client!")

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines
	go newClient.wsClientReader()
	go newClient.wsClientWriter()
}
