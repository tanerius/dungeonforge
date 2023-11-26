package main

import (
	"log"
	"net/http"

	"github.com/tanerius/dungeonforge/pkg/events"
	"github.com/tanerius/dungeonforge/pkg/server"
)

func main() {
	var eventManager *events.EventManager = events.NewEventManager()
	eventManager.RegisterEventType(events.EventClientConnected)    // event for when the client makes a websocket connection
	eventManager.RegisterEventType(events.EventClientDisconnected) // event for when client disconnects its socket connection
	eventManager.RegisterEventType(events.EventClientRegistered)   // event fires when client registered to hub
	eventManager.RegisterEventType(events.EventMessageReceived)    // event fires when client receives a message

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
