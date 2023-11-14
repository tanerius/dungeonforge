package main

import "github.com/tanerius/dungeonforge/pkg/server"

func main() {
	var gameServer *server.DungeonForge = server.NewGameServer()
	go gameServer.Run()

	var server *server.SocketServer = server.NewSocketServer(gameServer)
	server.StartHTTPServer()
	select {}
}
