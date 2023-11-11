package server

import (
	"time"

	log "github.com/sirupsen/logrus"

	gameLoop "github.com/kutase/go-gameloop"
)

type GameServer struct {
	gl        *gameLoop.GameLoop
	gameCoord *Coordinator
}

func NewGameServer(_c *Coordinator) *GameServer {
	return &GameServer{
		gl:        nil,
		gameCoord: _c,
	}
}

func (gs *GameServer) Run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Perform game logic here.
			// Update player data, calculate resources, etc.
			// Send game updates to connected players via their player.conn.
		/*

			case c := <-gs.gameCoord.register:
				// A new player has connected, you can perform any initialization here.
				log.Println("New player connected. " + c.entityId.String())

			case c := <-gs.gameCoord.unregister:
				// A player has disconnected, you can perform any cleanup here.
				log.Println("Player disconnected. " + c.entityId.String())
		*/
		case playerMessage := <-gs.gameCoord.playerMessages:
			// Handle player messages here.
			// Example: Log the received message.
			log.Printf("Received message from player: %s\n", playerMessage)
		}
	}
}
