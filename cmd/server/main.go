package main

import (
	"github.com/tanerius/dungeonforge/pkg/game"
	"github.com/tanerius/dungeonforge/pkg/server"
)

func main() {
	var gameServer *game.DungeonForge = game.NewDungeonForge()
	go gameServer.Run()

	var server *server.SocketServer = server.NewSocketServer(gameServer)
	server.StartHTTPServer()
	select {}
}
