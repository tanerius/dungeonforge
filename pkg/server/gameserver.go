package server

type GameServer interface {
	HandleClient(*client) error
	Run()
	Stop()
}
