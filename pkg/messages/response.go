package messages

// Game Specific data
type RspData interface{}

type Response struct {
	Ts   int64   `json:"ts"`
	Sid  string  `json:"sid"`
	Data RspData `json:"data,omitempty"`
}
