package game

import (
	"errors"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	game "github.com/tanerius/dungeonforge/pkg/game/messaes"
	"github.com/tanerius/dungeonforge/pkg/messages"
	"github.com/tanerius/dungeonforge/pkg/server"

	gameLoop "github.com/kutase/go-gameloop"
)

type DungeonForge struct {
	*server.GameConfig
	serverId        uuid.UUID
	gameloop        *gameLoop.GameLoop
	gameCoordinator *server.Coordinator
	isRunning       bool
}

func NewDungeonForge() *DungeonForge {

	return &DungeonForge{
		GameConfig: &server.GameConfig{
			GameId:      1,
			GameName:    "Dungeon Forge",
			IsTurnBased: false,
			IsRealtime:  false,
		},
		serverId:        uuid.New(),
		gameloop:        nil,
		gameCoordinator: server.NewCoordinator(),
		isRunning:       false,
	}
}

func (d *DungeonForge) Config() server.GameConfig {
	return server.GameConfig{
		GameId:      d.GameId,
		GameName:    d.GameName,
		IsTurnBased: d.IsTurnBased,
		IsRealtime:  d.IsRealtime,
	}
}

// A handler for new clients. Every new client should be handed off here!
func (d *DungeonForge) HandleClient(_client *server.Client) error {
	if !d.isRunning {
		return errors.New("gameserver not running")
	}

	go func() {
		d.gameCoordinator.Register <- _client
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
			connections, players := d.gameCoordinator.GetCounts()
			log.Printf("gameserver * %d players / %d connections", players, connections)
		case msg := <-d.gameCoordinator.PlayerMessagesChan:
			log.Printf("gameserver * received %v", msg)
			go d.ProcessMsg(msg)
		}
	}

}

// This is the place where messages from clients are processed in the game
func (d *DungeonForge) ProcessMsg(_msg *messages.Payload) {

	// for now just answer the client with an empty message
	resp := &messages.Response{
		Ts:  time.Now().Unix(),
		Sid: d.serverId.String(),
		Data: &game.Response{
			Ts:   time.Now().Unix(),
			Code: game.RspNotAuthorised,
			Msg:  "not authenticated",
		},
	}

	d.gameCoordinator.SendMessageToClient(resp, _msg.ClientId)
}

// Stop the gameserver
func (d *DungeonForge) Stop() {
	// TODO: Implement
	d.isRunning = true
}
