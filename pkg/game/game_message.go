package game

// Definition of game messages
type GameMessageEvent struct {
	ClientId string
	UserId   string
	Data     map[string]interface{}
}
