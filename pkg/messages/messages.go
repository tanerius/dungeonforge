package messages

type PlayerID string

const (
	CmdLost       int = iota // alient sends to server that peer was lost
	CmdExec                  // execute a game specific command
	CmdDisconnect            // execute a disconnect
)

type Request struct {
	Seq      int64 `json:"sqeuence,omitempty"`
	CmdType  int   `json:"cmdType,omitempty"`
	DataType int   `json:"dataType,omitempty"` // user defined type
}

type Response struct {
	Ts  int64  `json:"ts"`
	Sid string `json:"sid"`
	Cmd int    `json:"cmd"`
	Msg string `json:"msg"`
}
