package gameevents

import (
	"github.com/tanerius/dungeonforge/pkg/events"
	"github.com/tanerius/dungeonforge/pkg/messages"
)

const (
	GameEventLogin      events.EventIdType = "gLogin"
	GameEventDisconnect events.EventIdType = "gDisconnect"
)

const (
	TypeNothing int = iota // nothing
	TypeLogin              // Login player
	TypeLogout             // Logout player
)

type RequestLogin struct {
	messages.Request
	PlayerId string `json:"pid,omitempty"`
	Password string `json:"pass,omitempty"`
	ClientId string `json:"cid,omitempty"`
}

func (e *RequestLogin) EventId() events.EventIdType {
	return GameEventLogin
}

type RspCode int

const (
	RspOK            RspCode = 0
	RspNotAuthorised RspCode = 1
	RspError         RspCode = 3
)

type Response struct {
	Ts   int64       `json:"ts"`
	Code RspCode     `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
