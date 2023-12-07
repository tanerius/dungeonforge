package usermanagement

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/game"
	"github.com/tanerius/dungeonforge/pkg/server"
)

type _UserMessageHandler struct {
	registrar *Registrar
}

func NewUserMessageHandler(_r *Registrar) *_UserMessageHandler {
	return &_UserMessageHandler{
		registrar: _r,
	}
}

func (h *_UserMessageHandler) Type() int {
	return server.EventMsgReceived
}

func (h *_UserMessageHandler) HandleEvent(_event *eventgoround.Event) {
	log.Debugln("[_UserMessageHandler] handling event")
	msgEvent, err := eventgoround.GetEventData[*server.MessageEvent](_event)
	if err == nil {
		var jsonMap map[string]string
		err := json.Unmarshal(msgEvent.Data(), &jsonMap)

		if err != nil {
			log.Error(err)
		} else {
			val, ok := jsonMap["act"]
			if !ok {
				log.Errorf("malformed %s", msgEvent.ClientId())
				return
			} else {
				switch val {
				case "login":
					// player wants to login

				case "logout":
					// player want to log out providing the player is logged in and provides a token
					h.registrar.disconnectClient(msgEvent.ClientId())
				default:
					//any other case is a game message providing the player is logged in and provides a token
					if ok, _ := h.registrar.isValudUser(msgEvent.ClientId(), jsonMap); ok {
						gameMessage := &game.GameMessageEvent{
							ClientId: msgEvent.ClientId(),
							UserId:   h.registrar.clientToUser[msgEvent.ClientId()],
							Data:     jsonMap,
						}

						event := eventgoround.NewEvent(game.GameMsg, gameMessage)
						h.registrar.eventManager.DispatchPriorityEvent(event)
					}
				}
			}
		}

	} else {
		log.Error(err)
	}
}
