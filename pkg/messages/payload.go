package messages

type Command string

const (
	CmdPing    Command = "ping"    // ping
	CmdPong    Command = "pong"    // pong
	CmdGetChar Command = "getchar" // get character and coplayers
	CmdLvlUp   Command = "lvlup"   // level up a character trait
)

type CmdData interface{}

type Data1 struct {
	Id int `json:"id"`
}

type Payload struct {
	Token string  `json:"token"`
	Seq   int64   `json:"sqeuence"`
	Cmd   Command `json:"cmd"`
	Data  CmdData `json:"data,omitempty"`
}

var XData Payload = Payload{
	Token: "x",
	Seq:   33,
	Cmd:   CmdPing,
}

var YData Payload = Payload{
	Token: "y",
	Seq:   33,
	Cmd:   CmdLvlUp,
	Data: PersonJson{
		Name: "Tanerius",
		Age:  45,
	},
}
