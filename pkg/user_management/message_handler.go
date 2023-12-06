package usermanagement

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/server"
)

type _UserMessageHandler struct {
}

func NewUserMessageHandler() *_UserMessageHandler {
	return &_UserMessageHandler{}
}

func (h *_UserMessageHandler) Type() int {
	return server.EventMsgReceived
}

func (m *_UserMessageHandler) HandleEvent(_event *eventgoround.Event) {
	log.Debugln("[_UserMessageHandler] handling event")
	msgEvent, err := eventgoround.GetEventData[*server.MessageEvent](_event)
	if err == nil {
		myString := string(msgEvent.Data()[:])
		log.Debugf("[_UserMessageHandler] %s : %v", msgEvent.ClientId(), myString)
		var jsonMap map[string]interface{}
		err := json.Unmarshal(msgEvent.Data(), &jsonMap)

		if err != nil {
			log.Error(err)
		} else {
			fmt.Println(jsonMap)
		}

	} else {
		log.Error(err)
	}
}
