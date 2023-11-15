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
	players           map[PlayerID]uuid.UUID

	register           chan *client
	unregister         chan *client
	playerMessagesChan chan *messages.Payload
	playerLogoutChan   chan PlayerID

	playerLoginChan chan *struct {
		Pid PlayerID
		Cid uuid.UUID
	}
}

// Create a new Coordinator
func newCoordinator() *coordinator {
	return &coordinator{
		id:                uuid.New(), // handle this in case of panic
		activeConnections: make(clients),
		players:           make(map[PlayerID]uuid.UUID),

		register:           make(chan *client),
		unregister:         make(chan *client),
		playerMessagesChan: make(chan *messages.Payload, 32),
		playerLogoutChan:   make(chan PlayerID),
		playerLoginChan: make(chan *struct {
			Pid PlayerID
			Cid uuid.UUID
		}),
	}
}

// Run the Coordinator
func (hub *coordinator) Run() {
	log.Println("coordinator * started")
	for {
		select {
		//LOGOUT
		case playerLoggingOut := <-hub.playerLogoutChan:
			// make sure first that player logged out
			delete(hub.players, playerLoggingOut)

		//LOGIN
		case newPlayerConnection := <-hub.playerLoginChan:
			if currentConnectionId, ok := hub.players[newPlayerConnection.Pid]; ok {
				// player already registered
				// make sure his connection is the same as the new one registering
				if currentConnectionId != newPlayerConnection.Cid {
					// players new connection doesnt match his old...a new login is made
					// so disconnect old one
					log.Printf("Coordinator * Existing player %s connection changed %s -> %s  \n",
						newPlayerConnection.Pid, currentConnectionId.String(), newPlayerConnection.Cid.String())
					go func() {
						hub.unregister <- hub.activeConnections[currentConnectionId]
					}()
				} else {
					log.Warnf("Coordinator * Existing player retrying login...")
				}
			} else {
				//player hasnt been registered yet
				// first make sure his connection exists - should ALWAYS exist
				if _, ok := hub.activeConnections[newPlayerConnection.Cid]; ok {
					hub.players[newPlayerConnection.Pid] = newPlayerConnection.Cid
					log.Printf("Coordinator * New player to connection registered %s -> %s  \n", newPlayerConnection.Pid, newPlayerConnection.Cid.String())
				}
			}

		// REGISTER CONNECTION
		case c := <-hub.register:
			log.Printf("Coordinator * Client %s connected \n", c.clientId.String())
			hub.activeConnections[c.clientId] = c
			c.activateClientOnGameserver(hub)

		// UNREGISTER CONNECTION
		case c := <-hub.unregister:
			if _, ok := hub.activeConnections[c.clientId]; ok {
				log.Printf("Coordinator * Disconnecting %s \n", c.clientId.String())
				// TODO: make sure gameserver also receives word that a potential player is dropped

				delete(hub.activeConnections, c.clientId)
				// maybe we should also close the channel here
				c.cn.Close()

			}
		}

	}
}

// Get total connections and total logged in players
func (hub *coordinator) GetCounts() (int, int) {
	return len(hub.activeConnections), len(hub.players)
}

func (hub *coordinator) SendMessageToClient(_msg *messages.Response, _clients uuid.UUID) {
	hub.SendMessageToClients(_msg, []uuid.UUID{_clients})
}

func (hub *coordinator) SendMessageToClients(_msg *messages.Response, _clients []uuid.UUID) {
	for _, _client := range _clients {
		log.Printf("Coordinator * sending to %s \n", _client.String())
		if c, ok := hub.activeConnections[_client]; ok {
			go func() {
				log.Printf("Coordinator * seding... ")
				c.toSend <- _msg
				log.Printf("Coordinator * sent ")
			}()
			/*
				select {
				case c.toSend <- _msg:
					break
				default:
					c.closeRequested = true
					hub.unregister <- c
					log.Errorf("Coordinator * client %s broke \n", c.clientId.String())
				}
			*/
		}
	}
}

func (hub *coordinator) BroadcastMessageToClients(_msg *messages.Response) {

	for _, c := range hub.activeConnections {
		select {
		case c.toSend <- _msg:
			// Send successful
		default:
			c.closeRequested = true
			hub.unregister <- c
			log.Errorf("Coordinator * client %s broke \n", c.clientId.String())
		}
	}

}
