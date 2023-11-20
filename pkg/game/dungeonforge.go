package game

import (
	"errors"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
	players         map[messages.PlayerID]string
}

func NewDungeonForge() *DungeonForge {
	coordinator := server.NewCoordinator()
	go coordinator.Run()

	return &DungeonForge{
		GameConfig: &server.GameConfig{
			GameId:      1,
			GameName:    "Dungeon Forge",
			IsTurnBased: false,
			IsRealtime:  false,
		},
		serverId:        uuid.New(),
		gameloop:        nil,
		gameCoordinator: coordinator,
		isRunning:       true,
		players:         make(map[messages.PlayerID]string),
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
		return errors.New("[s] not running")
	}

	go d.processClient(_client)

	return nil
}

// Run a client
func (d *DungeonForge) processClient(_client *server.Client) {
	defer func() {
		_client.DeActivateClient()
		log.Infof("[s] deregistering %s ", _client.ID())
		d.gameCoordinator.Unregister <- _client
	}()

	// TODO: make a timeout here
	d.gameCoordinator.Register <- _client
	var writeChan chan *messages.Response = make(chan *messages.Response)
	var readChan chan *messages.Request = make(chan *messages.Request)

	if err := _client.ActivateClient(writeChan, readChan); err != nil {
		log.Error(err)
		return
	}

	for {
		msg, ok := <-readChan

		if !ok {
			log.Errorf("[s] can't read channel %s ", _client.ID())
			return
		}

		// make sure its not disconnect
		if msg.CmdType == messages.CmdLost {
			log.Debugf("[s] lost sonnection to  %s ", _client.ID())
			return
		}

		if msg.CmdType != messages.CmdExec {
			log.Errorf("[s] unknown message from %s ", _client.ID())
		} else {
			log.Debugf("[s] data %s : %v ", _client.ID(), msg)
		}
	}

	/*
		for {
			select {
			case msg := <-readChan:
				// Read messsage from the client
				log.Infof("Reading %v", msg)
				if msg.Cmd == messages.CmdDisconnect {
					log.Debugf("server received disconnect request from %s ", _client.ID())
					writeChan <- &messages.Response{
						Ts:  time.Now().Unix(),
						Sid: d.serverId.String(),
						Cmd: messages.CmdDisconnect,
						Msg: "bye",
					}
					return
				} else {
					writeChan <- d.responseNotAuthorized()
				}
			case <-_client.DisconnectedChan:
				// TODO: do all game disconnects here
				log.Debugf("client %s disconnected", _client.ID())
				return
			}
		}
	*/
}

// Stop the gameserver
func (d *DungeonForge) Stop() {
	// TODO: Implement
	d.isRunning = false
}
