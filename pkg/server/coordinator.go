package server

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/events"
)

// This is the connection coordinator responsible for keeping connection pools synced to a game server.
// A game coordinator controls player pools for a given server.
// This can be distributed.
// If coordinator dies, all users connected to italso disconnect from the game.
type Coordinator struct {
	id                string
	activeConnections clients
	eventManager      *events.EventManager
	register          chan *Client
	unregister        chan string
	broadcastChan     chan *MessageEvent
	sendToClient      chan *MessageEvent
	quitClient        chan string
}

// Create a new Coordinator
func NewCoordinator(_em *events.EventManager) *Coordinator {
	return &Coordinator{
		id:                uuid.NewString(), // handle this in case of panic
		activeConnections: make(clients),
		eventManager:      _em,
		register:          make(chan *Client),
		unregister:        make(chan string),
		broadcastChan:     make(chan *MessageEvent, 5),
		sendToClient:      make(chan *MessageEvent, 10),
		quitClient:        make(chan string, 3),
	}
}

// Events handler handler
func (hub *Coordinator) Handle(event events.Event) {
	switch resolvedEvent := event.(type) {
	case *ClientEvent:
		if resolvedEvent.EventId() == events.EventClientConnected {
			hub.register <- resolvedEvent.client
		} else if resolvedEvent.EventId() == events.EventClientDisconnected {
			hub.unregister <- resolvedEvent.clientId
		}
	default:
		log.Warnf("coordinator received an unhandled event %v %T", resolvedEvent, resolvedEvent)
	}
}

// broadcast to all clients blocks
func (hub *Coordinator) BroadcastToAllClients(_data []byte) {
	msg := &MessageEvent{
		data:    _data,
		eventId: events.EventNull,
	}
	hub.broadcastChan <- msg
}

// disconnect a client
func (hub *Coordinator) DisconnectClient(_clientId string) {
	hub.quitClient <- _clientId
}

// register to handle necessary messages
func (hub *Coordinator) RegisterHandlers() {
	// register connection and disonnection events
	err := hub.eventManager.RegisterHandler(events.EventClientConnected, hub)
	if err != nil {
		panic(err)
	}
	err = hub.eventManager.RegisterHandler(events.EventClientDisconnected, hub)
	if err != nil {
		panic(err)
	}
}

// Run the Coordinator
func (hub *Coordinator) Run() {
	log.Infoln("coordinator * started")

	for {
		select {
		// REGISTER CONNECTION
		case c := <-hub.register:
			log.Debugf("Coordinator * Client %s connected \n", c.clientId)
			hub.activeConnections[c.clientId] = c
			c.activateClient()
			hub.eventManager.Dispatch(NewClientEvent(events.EventClientRegistered, c.clientId, c))
		// UNREGISTER CONNECTION
		case c := <-hub.unregister:
			if _, ok := hub.activeConnections[c]; ok {
				log.Debugf("Coordinator * Disconnecting %s \n", c)
				delete(hub.activeConnections, c)
			}
		case message, ok := <-hub.broadcastChan:
			if !ok {
				log.Error("broadcast channel closed")
			} else {
				for _, client := range hub.activeConnections {
					select {
					case client.sendChannel <- message.data:
					default:
						log.Errorf("could not broadcast to %s \n", client.clientId)
						client.deActivateClient()
					}
				}
			}
		case message, ok := <-hub.sendToClient:
			if !ok {
				log.Error("sendToClient channel closed")
			} else {
				if client, ok := hub.activeConnections[message.clientId]; ok {
					select {
					case client.sendChannel <- message.data:
					default:
						log.Errorf("could not send to %s \n", client.clientId)
						client.deActivateClient()
					}
				}
			}
		case clientId, ok := <-hub.quitClient:
			if !ok {
				log.Error("quitClient channel closed")
			} else {
				if client, ok := hub.activeConnections[clientId]; ok {
					client.deActivateClient()
				}
			}
		}
	}
}

// sent to client blocks
func (hub *Coordinator) SendMessageToClient(_clientId string, _data []byte) {
	msg := &MessageEvent{
		clientId: _clientId,
		data:     _data,
		eventId:  events.EventNull,
	}
	hub.sendToClient <- msg
}

func (hub *Coordinator) RunsInOwnThread() bool {
	return false
}
