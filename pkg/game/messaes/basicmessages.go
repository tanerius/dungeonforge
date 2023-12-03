package gameevents

import (
	"github.com/tanerius/dungeonforge/pkg/messages"
)

const (
	TypeNothing int = 50 // nothing
	TypeLogin   int = 51 // Login player
	TypeLogout  int = 52 // Logout player
)

type RequestLogin struct {
	messages.Request
	PlayerId string `json:"pid,omitempty"`
	Password string `json:"pass,omitempty"`
	ClientId string `json:"cid,omitempty"`
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
