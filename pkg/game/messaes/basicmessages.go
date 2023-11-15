package game

type Command int

const (
	CmdNothing Command = 0 // nothing
	CmdLogin   Command = 1 // Login player
)

type Request struct {
	Cmd  Command     `json:"cmd"`
	Ts   int64       `json:"ts"`
	Data interface{} `json:"data,omitempty"`
}

type Response struct {
	IsOk bool        `json:"isOk"`
	Ts   int64       `json:"ts"`
	Data interface{} `json:"data,omitempty"`
}

type RequestLogin struct {
	PlayerId string `json:"pid"`
	Password string `json:"pass"`
}
