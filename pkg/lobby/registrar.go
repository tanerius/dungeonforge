package usermanagement

import (
	"encoding/json"
	"errors"
	"net/mail"

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

// Registrar loop
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

// Handle user login request
func (r *Registrar) login(_cid, _name, _pass string) {
	if (len(_name) < 4) || (len(_name) > 15) {
		iname := &entities.User{
			ResponseCode: entities.RespInvalidUser,
			ResponseMsg:  "invalid username",
		}

		if err := r.sendUserResponse(_cid, iname); err != nil {
			log.Errorln(err)
		}
		return
	}

	if usr, err := r.database.Login(_name, _pass); err != nil {
		usr = &entities.User{
			ResponseCode: entities.RespLoginError,
			ResponseMsg:  err.Error(),
		}
		if err := r.sendUserResponse(_cid, usr); err != nil {
			log.Errorln(err)
		}
	} else {
		usr.ClientId = _cid
		usr.ResponseCode = entities.RespOK
		usr.ResponseMsg = "login"

		r.clientToUser[_cid] = usr.ID.Hex()

		// check if user already has a connection
		_, userFound := r.onlineUsers[usr.ID.Hex()]
		if userFound {
			usr.ResponseMsg = "relogin"
			// relogin required
			r.onlineUsers[usr.ID.Hex()] = nil
		}

		r.clientToUser[_cid] = usr.ID.Hex()
		r.onlineUsers[usr.ID.Hex()] = usr
		if err := r.sendUserResponse(_cid, usr); err != nil {
			log.Errorln(err)
		}
	}
}

// Handle user logout request
func (r *Registrar) logout(_cid, token string) {
	if ok, _ := r.isValudUser(_cid, token); !ok {
		r.notAuthenticatedResponse(_cid)
		return
	}

	usr := &entities.User{
		ResponseCode: entities.RespOK,
		ResponseMsg:  "logout",
	}

	if err := r.sendUserResponse(_cid, usr); err != nil {
		log.Errorln(err)
	}

	r.disconnectClient(_cid)
}

// Handle register new user request
func (r *Registrar) register(_cid, _email, _pass, _name string) {
	// check email
	if !r.isValidEmail(_email) {
		ie := &entities.User{
			ResponseCode: entities.RespInvalidEmail,
			ResponseMsg:  "invalid email",
		}

		if err := r.sendUserResponse(_cid, ie); err != nil {
			log.Errorln(err)
		}
		return
	}

	// check username
	if (len(_name) < 4) || (len(_name) > 15) {
		iname := &entities.User{
			ResponseCode: entities.RespInvalidUser,
			ResponseMsg:  "invalid username",
		}

		if err := r.sendUserResponse(_cid, iname); err != nil {
			log.Errorln(err)
		}
		return
	}

	if usr, err := r.database.Register(_email, _pass, _name); err != nil {
		usr = &entities.User{
			ResponseCode: entities.RespRegisterError,
			ResponseMsg:  err.Error(),
		}
		if err := r.sendUserResponse(_cid, usr); err != nil {
			log.Errorln(err)
		}
	} else {
		log.Debugf("%v", usr)
		usr.ClientId = _cid
		usr.ResponseCode = entities.RespOK
		usr.ResponseMsg = "register"

		r.clientToUser[_cid] = usr.ID.Hex()
		r.onlineUsers[usr.ID.Hex()] = usr

		if err := r.sendUserResponse(_cid, usr); err != nil {
			log.Errorln(err)
		}
	}
}

// RESPONSES

// User isnt authenticated response. Bad token or invalid credentials
func (r *Registrar) notAuthenticatedResponse(_cid string) {
	r.customResponse(_cid, "not authenticated", entities.RespNotAuthenticated)
}

func (r *Registrar) malformedResponse(_cid string) {
	r.customResponse(_cid, "invalid request", entities.RespInvalidRequest)
}

func (r *Registrar) customResponse(_cid, _msg string, _code int) {
	usr := &entities.User{
		ResponseCode: _code,
		ResponseMsg:  _msg,
	}

	if err := r.sendUserResponse(_cid, usr); err != nil {
		log.Errorln(err)
	}
}

func (r *Registrar) sendUserResponse(_cid string, _usr *entities.User) error {
	b, err := json.Marshal(_usr)
	if err != nil {
		return err
	} else {
		go func() {
			r.coordinator.SendMessageToClient(_cid, b)
		}()
	}
	return nil
}

func (r *Registrar) isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
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
	userId, ok := r.clientToUser[_cid]
	if ok {
		// client is actually a user
		r.disconnectUser(userId)
		delete(r.clientToUser, _cid)
	}
}
