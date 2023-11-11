package main

import "github.com/tanerius/dungeonforge/pkg/server"

func main() {
	var coordinator *server.Coordinator = server.NewCoordinator()
	go coordinator.Run()

	var gameServer *server.GameServer = server.NewGameServer(coordinator)
	go gameServer.Run()

	var server *server.SocketServer = server.NewSocketServer(gameServer, coordinator)

	server.StartHTTPServer()
	select {}
}
