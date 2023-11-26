package main

import (
	"log"
	"net/http"

	"github.com/tanerius/dungeonforge/pkg/events"
	"github.com/tanerius/dungeonforge/pkg/game"
	"github.com/tanerius/dungeonforge/pkg/server"
)

func main() {
	var eventManager *events.EventManager = events.NewEventManager()
	eventManager.RegisterEventType(events.EventClientConnected)
	eventManager.RegisterEventType(events.EventClientDisconnected)
	go eventManager.Run()

	var gameServer *game.DungeonForge = game.NewDungeonForge()

	var server *server.Server = server.NewServer(gameServer, eventManager)
	go server.StartServer(true)

	log.Fatal(http.ListenAndServe(":40000", nil))
	eventManager.Stop()
}
