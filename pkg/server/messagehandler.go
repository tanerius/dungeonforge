package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/EventGoRound/eventgoround"
)

type StreamProcessor interface {
	ProcessData(string, []byte)
}

// This is a handle to determines what type of game message
type ClientMessageHandler struct {
	streamProcessor StreamProcessor
}

func NewClientMessageHandler(_sp StreamProcessor) *ClientMessageHandler {
	return &ClientMessageHandler{
		streamProcessor: _sp,
	}
}

func (h *ClientMessageHandler) Type() int {
	return EventMsgReceived
}

func (m *ClientMessageHandler) HandleEvent(_event *eventgoround.Event) {
	log.Debugln("[GameMessageHandler] handling event")
	msgEvent, err := eventgoround.GetEventData[*MessageEvent](_event)
	if err == nil {
		m.streamProcessor.ProcessData(msgEvent.ClientId(), msgEvent.Data())
		log.Debugln(msgEvent.Data())
	} else {
		log.Error(err)
	}
}
