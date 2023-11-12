package messages

type Command string

const (
	CmdPing    Command = "ping" // ping network specific
	CmdPong    Command = "pong" // pong network specific
	CmdGetChar Command = "exec" // execute a game specific command
)

// Game Specific data
type CmdData interface{}

type Payload struct {
	Token string  `json:"token"`
	Seq   int64   `json:"sqeuence"`
	Cmd   Command `json:"cmd"`
	Data  CmdData `json:"data,omitempty"`
}
