package usermanagement

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/database"
	"github.com/tanerius/dungeonforge/pkg/entities"
	"github.com/tanerius/dungeonforge/pkg/server"
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
	database     *database.DBWrapper
	coordinator  *server.Coordinator
}

func NewRegistrar(_em *eventgoround.EventManager, _db *database.DBWrapper, _c *server.Coordinator) *Registrar {
	return &Registrar{
		onlineUsers:  make(map[string]*entities.User),
		clientToUser: make(map[string]string),
		eventManager: _em,
		database:     _db,
		coordinator:  _c,
	}
}

func (r *Registrar) Run() {
	defer func() {
		// disconnect all users from database
	}()

	dcHandler := NewUserDisconnectHandler(r)
	r.eventManager.RegisterListener(dcHandler)

	userMessageHandler := NewUserMessageHandler(r)
	r.eventManager.RegisterListener(userMessageHandler)

	select {}
}

func (r *Registrar) logout(cid string, data map[string]string) {
	if ok, err := r.isValudUser(cid, data); !ok {
		log.Error(err)
		return
	}

	r.disconnectClient(cid)
}

func (r *Registrar) isValudUser(cid string, data map[string]string) (bool, error) {
	//TODO: validate token
	userId, ok := r.clientToUser[cid]
	if ok {
		user, userok := r.onlineUsers[userId]
		if !userok {
			return false, errors.New("user not online")
		}

		token, okToken := data["token"]
		if !okToken {
			return false, errors.New("invalid user token")
		}

		if user.Token != token {
			return false, errors.New("spoofed user")
		}

		return true, nil
	}

	return false, errors.New("spoofed peer")
}

func (r *Registrar) disconnectUser(uid string) {
	// do database disconnection here and remove from onlineUsers
	r.database.Logout(uid)
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
