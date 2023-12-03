package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/EventGoRound/eventgoround"
)

type DisconnectHandler struct {
	coordinator *Coordinator
}

func (h *DisconnectHandler) HandleEvent(_event *eventgoround.Event) {
	client, err := eventgoround.GetEventData[*ClientEvent](_event)
	if err != nil {
		log.Error(err)
	} else {

		go h.coordinator.DisconnectClient(client.clientId)
	}
}

func (h *DisconnectHandler) Type() int {
	return EventClientDisconnect
}

func NewDisconnectHandler(_c *Coordinator) *DisconnectHandler {
	return &DisconnectHandler{
		coordinator: _c,
	}
}

type ConnectHandler struct {
	coordinator *Coordinator
}

func (h *ConnectHandler) HandleEvent(_event *eventgoround.Event) {
	log.Debugln("[ConnectHandler] Handling connect event ...")
	client, err := eventgoround.GetEventData[*ClientEvent](_event)
	if err != nil {
		log.Error(err)
	} else {
		go h.coordinator.RegisterClient(client.client)
	}
}

func (h *ConnectHandler) Type() int {
	return EventClientConnect
}

func NewConnectHandler(_c *Coordinator) *ConnectHandler {
	return &ConnectHandler{
		coordinator: _c,
	}
}
