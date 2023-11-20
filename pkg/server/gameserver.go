package server

type GameServer interface {
	Config() GameConfig
	HandleClient(*Client) error
	Stop()
}
