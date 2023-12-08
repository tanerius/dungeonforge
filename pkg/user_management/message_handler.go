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
			token, hasToken := jsonMap["token"]

			val, ok := jsonMap["act"]
			if !ok {
				log.Errorf("malformed %s", msgEvent.ClientId())
				return
			} else {
				switch val {
				case "register":
					valEmail, okEmail := jsonMap["email"]
					valPass, okPass := jsonMap["pass"]

					if !(okEmail && okPass) {
						log.Errorf("malformed register %s", msgEvent.ClientId())
						return
					}

					h.registrar.register(msgEvent.ClientId(), valEmail, valPass)

				case "login":
					// player wants to login
					valEmail, okEmail := jsonMap["email"]
					valPass, okPass := jsonMap["pass"]

					if !(okEmail && okPass) {
						log.Errorf("malformed register %s", msgEvent.ClientId())
						return
					}

					h.registrar.login(msgEvent.ClientId(), valEmail, valPass)

				case "logout":

					if !hasToken {
						log.Errorf("malformed logout %s", msgEvent.ClientId())
						return
					}

					h.registrar.logout(msgEvent.ClientId(), token)
				default:
					if !hasToken {
						log.Errorf("malformed logout %s", msgEvent.ClientId())
						return
					}

					//any other case is a game message providing the player is logged in and provides a token
					if ok, _ := h.registrar.isValudUser(msgEvent.ClientId(), token); ok {
						uid := h.registrar.clientToUser[msgEvent.ClientId()]
						gameMessage := &game.GameMessageEvent{
							ClientId: msgEvent.ClientId(),
							User:     h.registrar.onlineUsers[uid],
							Data:     jsonMap,
						}

						event := eventgoround.NewEvent(game.GameMsg, gameMessage)
						go h.registrar.eventManager.DispatchPriorityEvent(event)
					}
				}
			}
		}

	} else {
		log.Error(err)
	}
}
