package game

type Command int

const (
	CmdNothing Command = 0 // nothing
	CmdLogin   Command = 1 // Login player
	CmdLogout  Command = 2 // Logout player
)

type Request struct {
	Cmd  Command     `json:"cmd"`
	Ts   int64       `json:"ts"`
	Data interface{} `json:"data,omitempty"`
}

type RequestLogin struct {
	PlayerId string `json:"pid"`
	Password string `json:"pass"`
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
