package lobby

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
				h.registrar.malformedResponse(msgEvent.ClientId())
				return
			} else {
				switch val {
				case "register":

					valEmail, okEmail := jsonMap["email"]
					valPass, okPass := jsonMap["pass"]
					valName, okName := jsonMap["name"]

					if !(okEmail && okPass && okName) {
						h.registrar.malformedResponse(msgEvent.ClientId())
						return
					}

					h.registrar.register(msgEvent.ClientId(), valEmail, valPass, valName)

					// send user to game loop processing
					h._SendMessage(msgEvent, nil)

				case "login":
					// player wants to login
					valName, okName := jsonMap["name"]
					valPass, okPass := jsonMap["pass"]

					if !(okName && okPass) {
						h.registrar.malformedResponse(msgEvent.ClientId())
						return
					}

					h.registrar.login(msgEvent.ClientId(), valName, valPass)

					// send user to game loop processing
					h._SendMessage(msgEvent, nil)

				case "logout":

					if !hasToken {
						h.registrar.notAuthenticatedResponse(msgEvent.ClientId())
						return
					}

					h.registrar.logout(msgEvent.ClientId(), token)
				default:
					if !hasToken {
						h.registrar.notAuthenticatedResponse(msgEvent.ClientId())
						return
					}

					//any other case is a game message providing the player is logged in and provides a token
					if ok, _ := h.registrar.isValudUser(msgEvent.ClientId(), token); ok {
						h._SendMessage(msgEvent, jsonMap)
					} else {
						h.registrar.notAuthenticatedResponse(msgEvent.ClientId())
					}
				}
			}
		}

	} else {
		log.Error(err)
	}
}

func (h *_UserMessageHandler) _SendMessage(_msgEvent *server.MessageEvent, _data map[string]string) {
	// send user to game loop processing
	uid := h.registrar.clientToUser[_msgEvent.ClientId()]
	gameMessage := &game.GameMessageEvent{
		ClientId: _msgEvent.ClientId(),
		User:     h.registrar.onlineUsers[uid],
		Data:     _data,
	}

	event := eventgoround.NewEvent(game.GameMsg, gameMessage)
	go h.registrar.eventManager.DispatchPriorityEvent(event)
}
