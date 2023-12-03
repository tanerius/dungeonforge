package main

import (
	"log"
	"net/http"

	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/server"
)

func main() {
	var eventManager *eventgoround.EventManager = eventgoround.NewEventManager()

	var coordinator *server.Coordinator = server.NewCoordinator(eventManager)
	coordinator.RegisterHandlers()
	go coordinator.Run()

	//var gameServer *game.DungeonForge = game.NewDungeonForge(coordinator, eventManager)

	var server *server.SocketServer = server.NewSocketServer(eventManager)
	go server.StartServer(true)

	// Run this last to prevent case where someone registers an event after event loop is running
	go eventManager.Run()

	log.Fatal(http.ListenAndServe(":40000", nil))
	eventManager.Stop()
}
