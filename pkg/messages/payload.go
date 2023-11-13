package messages

import "github.com/google/uuid"

type Command string

const (
	CmdPing     Command = "ping"     // ping network specific
	CmdPong     Command = "pong"     // pong network specific
	CmdValidate Command = "validate" // special command used to validate a connection and map playerID to clientID seq=1
	CmdGetChar  Command = "exec"     // execute a game specific command
)

// Game Specific data
type CmdData interface{}

type Payload struct {
	ClientId uuid.UUID
	Token    string  `json:"token"`
	Seq      int64   `json:"sqeuence"`
	Cmd      Command `json:"cmd"`
	Data     CmdData `json:"data,omitempty"`
}
