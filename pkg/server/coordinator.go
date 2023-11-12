package server

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/messages"
)

// This is the connection coordinator responsible for keeping connection pools synced to a game server.
// A game coordinator controls player pools for a given server.
// This can be distributed.
// If coordinator dies, all users connected to italso disconnect from the game.
type Coordinator struct {
	id                uuid.UUID
	activeConnections clients
	register          chan *client
	unregister        chan *client
	playerMessages    chan *messages.Payload
	gameServer        *GameServer
}

// Create a new Coordinator
func NewCoordinator(_forGameServer *GameServer) *Coordinator {
	return &Coordinator{
		id:                uuid.New(), // handle this in case of panic
		activeConnections: make(clients),
		register:          make(chan *client),
		unregister:        make(chan *client),
		playerMessages:    make(chan *messages.Payload),
		gameServer:        _forGameServer,
	}
}

// Run the Coordinator
func (hub *Coordinator) Run() {
	for {
		select {
		case c := <-hub.register:
			log.Printf("Coordinator * Client %s connected \n", c.clientId.String())
			hub.activeConnections[c.clientId] = c
			c.activateClientOnGameserver()
		case c := <-hub.unregister:
			if _, ok := hub.activeConnections[c.clientId]; ok {
				log.Printf("Coordinator * Disconnecting %s \n", c.clientId.String())
				delete(hub.activeConnections, c.clientId)
				c.cn.Close()
			}
		}
	}
}
