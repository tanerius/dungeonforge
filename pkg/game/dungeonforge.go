package game

import (
	"github.com/google/uuid"
	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/messages"
	"github.com/tanerius/dungeonforge/pkg/server"
)

type DungeonForge struct {
	*server.GameConfig
	serverId        string
	gameCoordinator *server.Coordinator
	eventManager    *eventgoround.EventManager
	isRunning       bool
	players         map[messages.PlayerID]*player
}

func NewDungeonForge(_hub *server.Coordinator, _em *eventgoround.EventManager) *DungeonForge {
	return &DungeonForge{
		GameConfig: &server.GameConfig{
			GameId:      1,
			GameName:    "Dungeon Forge",
			IsTurnBased: false,
			IsRealtime:  false,
		},
		serverId:        uuid.NewString(),
		eventManager:    _em,
		gameCoordinator: _hub,
		isRunning:       true,
		players:         make(map[messages.PlayerID]*player),
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

/*
// A handler for filtering message types
func (d *DungeonForge) Handle(event events.Event) {
	switch resolvedEvent := event.(type) {
	case *gameevents.RequestLogin:
		log.Infof("game login event %v %T", resolvedEvent, resolvedEvent)
		resp := &messages.Response{
			Ts:      time.Now().Unix(),
			Tokenid: "testToken",
			Sid:     d.serverId,
		}

		if data, err := json.Marshal(resp); err != nil {
			log.Error(err)
		} else {
			go d.gameCoordinator.SendMessageToClient(resolvedEvent.ClientId, data)
			return
		}
		return
	case *server.MessageEvent:
		if resolvedEvent.EventId() == gameevents.GameEventDisconnect {
			d.gameCoordinator.DisconnectClient(resolvedEvent.ClientId())
			return
		} else {
			log.Warnf("game received an unhandled event %v %T", resolvedEvent, resolvedEvent)
		}
	default:
		log.Warnf("game received an unhandled event %v %T", resolvedEvent, resolvedEvent)
	}
}
*/

// register to all events relevant for the game
func (d *DungeonForge) RegisterHandlers() {
	d.eventManager.RegisterListener(NewGameMessageHandler(d))
}

func (d *DungeonForge) RunsInOwnThread() bool {
	return false
}
