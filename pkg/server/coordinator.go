package server

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/EventGoRound/eventgoround"
)

// This is the connection coordinator responsible for keeping connection pools synced to a game server.
// A game coordinator controls player pools for a given server.
// This can be distributed.
// If coordinator dies, all users connected to italso disconnect from the game.
type Coordinator struct {
	id                string
	activeConnections clients
	eventManager      *eventgoround.EventManager
	register          chan *Client
	unregister        chan string
	broadcastChan     chan *MessageEvent
	sendToClient      chan *MessageEvent
	quitClient        chan string
}

// Create a new Coordinator
func NewCoordinator(_em *eventgoround.EventManager) *Coordinator {
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

func (hub *Coordinator) RegisterClient(_client *Client) {
	hub.register <- _client
}

func (hub *Coordinator) unregisterClient(_clientId string) {
	hub.unregister <- _clientId
}

// broadcast to all clients blocks
func (hub *Coordinator) BroadcastToAllClients(_data []byte) {
	msg := &MessageEvent{
		data: _data,
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
	hub.eventManager.RegisterListener(NewConnectHandler(hub))
	hub.eventManager.RegisterListener(NewDisconnectHandler(hub))
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
			event := eventgoround.NewEvent(EventClientRegistered, NewClientEvent(c.clientId, c))
			c.eventManager.DispatchPriorityEvent(event)
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
	}
	hub.sendToClient <- msg
}

func (hub *Coordinator) RunsInOwnThread() bool {
	return false
}
