package main

import (
	"errors"
	"net"

	"github.com/tanerius/dungeonforge/pkg/logging"
	"github.com/tanerius/dungeonforge/src/lobby"
	"google.golang.org/grpc"
)

func main() {
	log := logging.NewLogger()
	log.LogInfo("Starting a lobby")
	service := lobby.NewMockedLobby(log)

	if service == nil {
		log.LogError(errors.New("no lobby service"), "can't create lobby service")
		return
	}

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.LogError(err, "")
		return
	}

	serverRegistrar := grpc.NewServer()
	lobby.RegisterLobbyServer(serverRegistrar, service)

	serverError := serverRegistrar.Serve(listener)

	if serverError != nil {
		log.LogError(serverError, "")
		return
	}
}
