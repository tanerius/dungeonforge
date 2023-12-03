package gameevents

import (
	"github.com/tanerius/dungeonforge/pkg/messages"
)

// these are all he events that can happen in the game
// starting them from 50 as the first 50 might be systemmessages
const (
	GameEventLogin  int = 50 // Login player
	GameEventLogout int = 51 // Logout player
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
