package game

import (
	"time"

	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/database"
	"github.com/tanerius/dungeonforge/pkg/server"
)

type GameServer struct {
	playerChannels map[string]chan *GameMessageEvent
	hub            *server.Coordinator
	db             *database.DBWrapper
	em             *eventgoround.EventManager
}

func NewGameServer(_h *server.Coordinator) *GameServer {
	return &GameServer{
		playerChannels: make(map[string]chan *GameMessageEvent),
		hub:            _h,
	}
}

func (g *GameServer) Run() {
	ticker := time.NewTicker(10 * time.Second)
}
