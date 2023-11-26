package server

type GameServerss interface {
	Config() GameConfig
	HandleClient(*Client) error
	Stop()
}
