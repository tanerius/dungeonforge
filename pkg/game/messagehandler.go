package game

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/events"
	game "github.com/tanerius/dungeonforge/pkg/game/messaes"
	gameevents "github.com/tanerius/dungeonforge/pkg/game/messaes"
	"github.com/tanerius/dungeonforge/pkg/messages"
	"github.com/tanerius/dungeonforge/pkg/server"
)

type GameMessageHandler struct {
	eventManager *events.EventManager
}

func NewMessageHandler(_em *events.EventManager) *GameMessageHandler {
	return &GameMessageHandler{
		eventManager: _em,
	}
}

func (m *GameMessageHandler) RegisterEvents() {
	err := m.eventManager.RegisterHandler(events.EventMessageReceived, m)
	if err != nil {
		panic(err)
	}

}

func (m *GameMessageHandler) Handle(event events.Event) {
	log.Debugln("GOT MESSAGE")
	switch resolvedEvent := event.(type) {
	case *server.MessageEvent:
		var msg *messages.Request = &messages.Request{}
		clientId := resolvedEvent.ClientId()
		//var msg *game.RequestLogin = &game.RequestLogin{}

		if err := json.Unmarshal(resolvedEvent.Data(), msg); err != nil {
			log.Errorf("[s] cannot unmarshal message from %s : %v", clientId, err)
			return
		} else {
			log.Debugf("[s] data %s : %v ", clientId, msg)

			if msg.CmdType == messages.CmdDisconnect {
				// client requested disconnect
				m.eventManager.Dispatch(server.NewMessageEvent(gameevents.GameEventLogin, clientId, nil))
				return
			} else if msg.CmdType == messages.CmdExec {
				switch msg.DataType {
				case int(game.TypeLogin):
					var loginInfo *game.RequestLogin = &game.RequestLogin{}
					if err := json.Unmarshal(resolvedEvent.Data(), loginInfo); err != nil {
						log.Error(err)
						return
					} else {
						loginInfo.ClientId = clientId
						m.eventManager.Dispatch(loginInfo)
					}
				default:
					log.Errorf("unhandled game message type %v \n", msg)
				}
			}
		}
	default:
		log.Warnf("messageHandler received an unhandled event %v %T", resolvedEvent, resolvedEvent)
	}
}

func (m *GameMessageHandler) RunsInOwnThread() bool {
	return true
}
