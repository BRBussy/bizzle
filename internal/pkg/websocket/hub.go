package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
)

type Hub struct {
	/*
	   A central hub will receive all incoming messages and broadcast them
	   to all registered "Client"s
	   (i.e. the Client structures in the clients map)
	*/
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client

	content []byte
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	/*
	   Channels are like FIFO Stacks
	   A Client will store a message in one of the 3 channels
	   A go routine will unstack them as soon as possible by arrival
	   time.
	*/
	log.Info().Msg("Starting websocket connection hub")
	for {
		select {
		case c := <-h.Register:
			log.Info().Msg("Client Registered to hub")
			h.Clients[c] = true
			helloClient, err := json.Marshal(message{
				Type: "WSS",
				Data: "Welcome to the hub :)",
			})
			if err == nil {
				c.Send <- helloClient
			} else {
				log.Warn().Msg("Unable to marshal message for client")
			}
			break

		case c := <-h.Unregister:
			_, ok := h.Clients[c]
			if ok {
				delete(h.Clients, c)
				close(c.Send)
			}
			break

		case m := <-h.Broadcast:
			fmt.Println("Broadcast!")
			h.content = m
			h.broadcastMessage()
			break
		}
	}
}

func (h *Hub) broadcastMessage() {
	for c := range h.Clients {
		select {
		case c.Send <- h.content:
			break

		// We can't reach the client
		default:
			close(c.Send)
			delete(h.Clients, c)
		}
	}
}
