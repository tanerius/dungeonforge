package usermanagement

import (
	"encoding/json"
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

func (r *Registrar) login(_cid, _email, _pass string) {
	if usr, err := r.database.Login(_email, _pass); err != nil {
		log.Errorln(err)
	} else {
		usr.ClientId = _cid
		usr.ResponseCode = entities.RespOK
		usr.ResponseMsg = ""

		r.clientToUser[_cid] = usr.ID.Hex()

		// check if user already has a connection
		existingUser, userFound := r.onlineUsers[usr.ID.Hex()]
		if userFound {
			usr.ResponseMsg = "relogin"
			exId := existingUser.ClientId
			// relogin required
			r.onlineUsers[usr.ID.Hex()] = nil
			if exId != usr.ClientId {
				go r.coordinator.DisconnectClient(exId)
			}
		}

		log.Debugf("%v", usr)
		r.clientToUser[_cid] = usr.ID.Hex()
		r.onlineUsers[usr.ID.Hex()] = usr
		b, err := json.Marshal(usr)
		if err != nil {
			log.Error(err)
		} else {
			go func() {
				r.coordinator.SendMessageToClient(_cid, b)
			}()
		}
	}
}

func (r *Registrar) logout(cid, token string) {
	if ok, err := r.isValudUser(cid, token); !ok {
		log.Error(err)
		return
	}
	go r.coordinator.DisconnectClient(cid)
	r.disconnectClient(cid)
}

func (r *Registrar) register(_cid, _email, _pass string) {
	if usr, err := r.database.Register(_email, _pass); err != nil {
		log.Errorln(err)
	} else {
		log.Debugf("%v", usr)
		usr.ClientId = _cid
		usr.ResponseCode = entities.RespOK
		usr.ResponseMsg = ""

		r.clientToUser[_cid] = usr.ID.Hex()
		r.onlineUsers[usr.ID.Hex()] = usr
		b, err := json.Marshal(usr)
		if err != nil {
			log.Errorln(err)
		} else {
			go func() {
				r.coordinator.SendMessageToClient(_cid, b)
			}()
		}
	}
}

func (r *Registrar) isValudUser(cid string, token string) (bool, error) {
	userId, ok := r.clientToUser[cid]
	if ok {
		user, userok := r.onlineUsers[userId]
		if !userok {
			return false, errors.New("user not online")
		}

		if user == nil {
			return false, errors.New("no user")
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

func (r *Registrar) disconnectClient(_cid string) {
	go r.coordinator.DisconnectClient(_cid)
	userId, ok := r.clientToUser[_cid]
	if ok {
		// client is actually a user
		r.disconnectUser(userId)
		delete(r.clientToUser, _cid)
	}
}
