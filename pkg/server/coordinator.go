package server

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/messages"
)

type PlayerID string

// This is the connection coordinator responsible for keeping connection pools synced to a game server.
// A game coordinator controls player pools for a given server.
// This can be distributed.
// If coordinator dies, all users connected to italso disconnect from the game.
type coordinator struct {
	id                uuid.UUID
	activeConnections clients
	register          chan *client
	unregister        chan *client
	playerMessage     chan *messages.Payload
}

// Create a new Coordinator
func newCoordinator() *coordinator {
	return &coordinator{
		id:                uuid.New(), // handle this in case of panic
		activeConnections: make(clients),
		register:          make(chan *client),
		unregister:        make(chan *client),
		playerMessage:     make(chan *messages.Payload),
	}
}

// Run the Coordinator
func (hub *coordinator) Run() {
	log.Println("coordinator * started")
	for {
		select {
		case c := <-hub.register:
			log.Printf("Coordinator * Client %s connected \n", c.clientId.String())
			hub.activeConnections[c.clientId] = c
			c.activateClientOnGameserver(hub)
		case c := <-hub.unregister:
			if _, ok := hub.activeConnections[c.clientId]; ok {
				log.Printf("Coordinator * Disconnecting %s \n", c.clientId.String())
				delete(hub.activeConnections, c.clientId)
				c.cn.Close()
			}
		}

	}
}
