package server

import "github.com/tanerius/dungeonforge/pkg/messages"

type GameServer interface {
	Config() GameConfig
	HandleClient(*Client) error
	ProcessMsg(*messages.Payload)
	Run()
	Stop()
}
