package server

import (
	"errors"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/messages"

	gameLoop "github.com/kutase/go-gameloop"
)

type DungeonForge struct {
	id              uuid.UUID
	gameloop        *gameLoop.GameLoop
	gameCoordinator *coordinator
	isRunning       bool
	messagePump     chan *messages.Payload
	players         map[PlayerID]struct {
		cid  uuid.UUID
		gcid uuid.UUID
	}
}

func NewGameServer() *DungeonForge {

	return &DungeonForge{
		id:              uuid.New(),
		gameloop:        nil,
		gameCoordinator: newCoordinator(),
		isRunning:       false,
		messagePump:     make(chan *messages.Payload, 32),
	}
}

// A handler for new clients. Every new client should be handed off here!
func (d *DungeonForge) HandleClient(_client *client) error {
	if !d.isRunning {
		return errors.New("gameserver not running")
	}

	go func() {
		d.gameCoordinator.register <- _client
	}()

	return nil
}

// Run the gameserver
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
			log.Printf("gameserver * ticking. %d players online", len(d.players))
		case msg := <-d.messagePump:
			log.Printf("gameserver * received %v", msg)
			go d.ProcessMsg(msg)
		}
	}

}

// This is the place where messages from clients are processed in the game
func (d *DungeonForge) ProcessMsg(_msg *messages.Payload) {
	// for now just anser the client with an empty message
	resp := &messages.Response{
		Ts:  time.Now().Unix(),
		Sid: d.id.String(),
	}

	d.gameCoordinator.SendMessageToClient(resp, _msg.ClientId)
}

// Stop the gameserver
func (d *DungeonForge) Stop() {
	// TODO: Implement
	d.isRunning = true
}
