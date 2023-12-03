package game

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/messages"
	"github.com/tanerius/dungeonforge/pkg/server"
)

// This is a handle to determines what type of game message
type GameMessageHandler struct {
	game *DungeonForge
}

func NewGameMessageHandler(_g *DungeonForge) *GameMessageHandler {
	return &GameMessageHandler{
		game: _g,
	}
}

func (h *GameMessageHandler) Type() int {
	return server.EventMsgReceived
}

func (m *GameMessageHandler) HandleEvent(_event *eventgoround.Event) {
	log.Debugln("[GameMessageHandler] handling event")
	msgEvent, err := eventgoround.GetEventData[*server.MessageEvent](_event)
	if err == nil {
		var msg *messages.Request = &messages.Request{}
		clientId := msgEvent.ClientId()

		if err := json.Unmarshal(msgEvent.Data(), msg); err != nil {
			log.Errorf("[GameMessageHandler] cannot unmarshal message from %s : %v", clientId, err)
			return
		} else {
			log.Debugf("[s] data %s : %v ", clientId, msg)

			if msg.CmdType == messages.CmdDisconnect {
				// client requested disconnect
				log.Debug("DISPATCH DISCONNECT EVENT")
				return
			} else if msg.CmdType == messages.CmdExec {
				log.Debug("DISPATCH EXEC EVENT")
			}
		}

	} else {
		log.Error(err)
	}
}
