package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/messages"
)

// This is the connection coordinator responsible for leeping connection pools synced
type Coordinator struct {
	activeConnections connections
	register          chan *connection
	unregister        chan *connection
	playerMessages    chan *messages.Payload
}

// Create a new Coordinator
func NewCoordinator() *Coordinator {
	return &Coordinator{
		activeConnections: make(connections),
		register:          make(chan *connection),
		unregister:        make(chan *connection),
		playerMessages:    make(chan *messages.Payload),
	}
}

// Run the Coordinator
func (hub *Coordinator) Run() {
	for {
		select {
		case c := <-hub.register:
			log.Printf("Coordinator * Client %s connected \n", c.entityId.String())
			hub.activeConnections[c.entityId] = c
		case c := <-hub.unregister:
			if _, ok := hub.activeConnections[c.entityId]; ok {
				log.Printf("Coordinator * Disconnecting %s \n", c.entityId.String())
				delete(hub.activeConnections, c.entityId)
				c.cn.Close()
			}
		}
	}
}
