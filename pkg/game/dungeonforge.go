package game

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
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
	serverId        string
	gameloop        *gameLoop.GameLoop
	gameCoordinator *server.Coordinator
	isRunning       bool
	players         map[messages.PlayerID]string
	mu              sync.Mutex
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
		serverId:        uuid.NewString(),
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

	d.mu.Lock()
	serverId := d.serverId
	d.mu.Unlock()

	token := uuid.NewString()

	defer func() {
		_client.DeActivateClient()
		log.Infof("[s] deregistering %s ", _client.ID())
		d.gameCoordinator.Unregister <- _client
	}()

	// TODO: make a timeout here
	d.gameCoordinator.Register <- _client
	var writeChan chan []byte = make(chan []byte)
	var readChan chan []byte = make(chan []byte)

	if err := _client.ActivateClient(writeChan, readChan); err != nil {
		log.Error(err)
		return
	}

	for {
		stream, ok := <-readChan

		if !ok {
			log.Errorf("[s] can't read channel %s ", _client.ID())
			return
		}

		var msg *messages.Request = &messages.Request{}
		//var msg *game.RequestLogin = &game.RequestLogin{}

		if err := json.Unmarshal(stream, msg); err != nil {
			log.Errorf("[s] cannot unmarshal message from %s : %v", _client.ID(), err)
		} else {
			log.Debugf("[s] data %s : %v ", _client.ID(), msg)

			if msg.CmdType == messages.CmdDisconnect {
				// disconnecting client
				close(writeChan)
				return
			} else if msg.CmdType == messages.CmdExec {
				ctx := context.Background()

				switch msg.DataType {
				case int(game.TypeLogin):
					if err := d.processLogin(ctx, serverId, token, stream, writeChan); err != nil {
						log.Errorf("[s] %v ", err)
					}
				default:
					log.Debugf("[s] data %v ", msg)
				}
			}
		}
	}
}

// Stop the gameserver
func (d *DungeonForge) Stop() {
	// TODO: Implement
	d.isRunning = false
}

func (d *DungeonForge) processLogin(ctx context.Context, sid string, token string, data []byte, writer chan<- []byte) error {
	// Create a timeout for the operation
	timeout := time.After(1 * time.Second)

	var loginInfo *game.RequestLogin = &game.RequestLogin{}

	if err := json.Unmarshal(data, loginInfo); err != nil {
		return err
	} else {
		log.Debugf("[s] login data %v ", loginInfo)
		resp := &messages.Response{
			Ts:      time.Now().Unix(),
			Tokenid: token,
			Sid:     sid,
		}

		if data, err := json.Marshal(resp); err != nil {
			return err
		} else {
			select {
			case writer <- data:
				return nil
			case <-timeout:
				return errors.New("processLogin timeout")
			}
		}
	}
}
