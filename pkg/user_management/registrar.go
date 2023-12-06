package usermanagement

import (
	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/entities"
)

const (
	UserAuthenticatedEvent int = 50
	UserLoggedOutEvent     int = 51
)

// Structure holding information on current users in the system
type Registrar struct {
	onlineUsers  map[string]*entities.User
	clientToUser map[string]string
	eventManager *eventgoround.EventManager
}

func NewRegistrar(_em *eventgoround.EventManager) *Registrar {
	return &Registrar{
		onlineUsers:  make(map[string]*entities.User),
		clientToUser: make(map[string]string),
		eventManager: _em,
	}
}

func (r *Registrar) Run() {
	defer func() {
		// disconnect all users from database
	}()

	dcHandler := NewUserDisconnectHandler(r)
	r.eventManager.RegisterListener(dcHandler)

	userMessageHandler := NewUserMessageHandler()
	r.eventManager.RegisterListener(userMessageHandler)

	select {}
}

func (r *Registrar) disconnectUser(uid string) {
	// do database disconnection here and remove from onlineUsers
	// TODO: probably a good idea to make this thread safe
	delete(r.onlineUsers, uid)
}

func (r *Registrar) disconnectClient(cid string) {
	userId, ok := r.clientToUser[cid]
	if ok {
		// client is actually a user
		r.disconnectUser(userId)
		delete(r.clientToUser, cid)
	}
}
