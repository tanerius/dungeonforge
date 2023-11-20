package messages

type Command int
type PlayerID string

const (
	CmdLost       Command = 0 // alient sends to server that peer was lost
	CmdExec       Command = 1 // execute a game specific command
	CmdDisconnect Command = 2 // execute a disconnect
)

type Request struct {
	ClientId string
	Seq      int64   `json:"sqeuence"`
	CmdType  Command `json:"cmdType"`
}

type Response struct {
	Ts  int64   `json:"ts"`
	Sid string  `json:"sid"`
	Cmd Command `json:"cmd"`
	Msg string  `json:"msg"`
}
