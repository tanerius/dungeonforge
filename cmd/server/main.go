package main

import (
	"log"
	"net/http"
	"time"

	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/database"
	"github.com/tanerius/dungeonforge/pkg/game"
	lobby "github.com/tanerius/dungeonforge/pkg/lobby"
	"github.com/tanerius/dungeonforge/pkg/server"
)

func main() {
	var eventManager *eventgoround.EventManager = eventgoround.NewEventManager()
	var lobbyDb *database.DBWrapper = database.NewDatabase()

	var coordinator *server.Coordinator = server.NewCoordinator(eventManager)
	coordinator.RegisterHandlers()
	go coordinator.Run()

	var server *server.SocketServer = server.NewSocketServer(eventManager)
	go server.StartServer(true)

	var registrar *lobby.Registrar = lobby.NewRegistrar(eventManager, lobbyDb, coordinator)
	go registrar.Run()

	var gameServer *game.GameServer = game.NewGameServer(coordinator, eventManager)
	go gameServer.Run()

	time.Sleep(1 * time.Second)
	// Run this last to prevent case where someone registers an event after event loop is running
	go eventManager.Run()

	log.Fatal(http.ListenAndServe(":40000", nil))
	eventManager.Stop()
}
