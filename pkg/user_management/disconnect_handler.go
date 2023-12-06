package usermanagement

import (
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/server"
)

// This is a handle to determines what type of game message
type UserDisconnectHandler struct {
	registrar *Registrar
}

func NewUserDisconnectHandler(_registrar *Registrar) *UserDisconnectHandler {
	return &UserDisconnectHandler{
		registrar: _registrar,
	}
}

func (h *UserDisconnectHandler) Type() int {
	return server.EventClientDisconnect
}

func (m *UserDisconnectHandler) HandleEvent(_event *eventgoround.Event) {
	log.Debugln("[UserDisconnectHandler] handling event")
	msgEvent, err := eventgoround.GetEventData[*server.ClientEvent](_event)
	if err == nil {
		userId, ok := m.registrar.clientToUser[msgEvent.ClientId()]
		if ok {
			// client is actually a user
			m.registrar.disconnectUser(userId)
			delete(m.registrar.clientToUser, msgEvent.ClientId())
		}

	} else {
		log.Error(err)
	}
}
