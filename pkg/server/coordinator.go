package server

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// This is the connection coordinator responsible for keeping connection pools synced to a game server.
// A game coordinator controls player pools for a given server.
// This can be distributed.
// If coordinator dies, all users connected to italso disconnect from the game.
type Coordinator struct {
	id                string
	activeConnections clients

	Register   chan *Client
	Unregister chan *Client
}

// Create a new Coordinator
func NewCoordinator() *Coordinator {
	return &Coordinator{
		id:                uuid.NewString(), // handle this in case of panic
		activeConnections: make(clients),

		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run the Coordinator
func (hub *Coordinator) Run() {
	log.Infoln("coordinator * started")
	for {
		select {
		// REGISTER CONNECTION
		case c := <-hub.Register:
			log.Debugf("Coordinator * Client %s connected \n", c.clientId)
			hub.activeConnections[c.clientId] = c

		// UNREGISTER CONNECTION
		case c := <-hub.Unregister:
			if _, ok := hub.activeConnections[c.clientId]; ok {
				log.Debugf("Coordinator * Disconnecting %s \n", c.clientId)
				delete(hub.activeConnections, c.clientId)
			}
		}

	}
}

// Get total connections and total logged in players
func (hub *Coordinator) GetCounts() int {
	return len(hub.activeConnections)
}
