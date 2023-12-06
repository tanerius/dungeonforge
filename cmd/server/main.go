package main

import (
	"log"
	"net/http"
	"time"

	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/server"
	usermanagement "github.com/tanerius/dungeonforge/pkg/user_management"
)

func main() {
	var eventManager *eventgoround.EventManager = eventgoround.NewEventManager()

	var coordinator *server.Coordinator = server.NewCoordinator(eventManager)
	coordinator.RegisterHandlers()
	go coordinator.Run()

	var server *server.SocketServer = server.NewSocketServer(eventManager)
	go server.StartServer(true)

	var registrar *usermanagement.Registrar = usermanagement.NewRegistrar(eventManager)
	go registrar.Run()

	time.Sleep(1 * time.Second)
	// Run this last to prevent case where someone registers an event after event loop is running
	go eventManager.Run()

	log.Fatal(http.ListenAndServe(":40000", nil))
	eventManager.Stop()
}
