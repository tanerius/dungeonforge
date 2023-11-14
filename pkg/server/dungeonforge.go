package server

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	gameLoop "github.com/kutase/go-gameloop"
)

type DungeonForge struct {
	gameloop        *gameLoop.GameLoop
	gameCoordinator *coordinator
	isRunning       bool
}

func NewGameServer() *DungeonForge {
	return &DungeonForge{
		gameloop:        nil,
		gameCoordinator: newCoordinator(),
		isRunning:       false,
	}
}

func (d *DungeonForge) HandleClient(_client *client) error {
	if !d.isRunning {
		return errors.New("gameserver not running")
	}

	go func() {
		d.gameCoordinator.register <- _client
	}()

	return nil
}

func (d *DungeonForge) Run() {
	log.Println("gameserver * starting...")
	if d.isRunning {
		log.Println("gameserver * alredy running")
		return
	}
	log.Println("gameserver * starting coordinator...")
	go d.gameCoordinator.Run()
	d.isRunning = true
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	log.Println("gameserver * going into main loop...")
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
			log.Println("gameserver * ticking...")

		}
	}
}

func (d *DungeonForge) Stop() {
	// TODO: Implement
}
