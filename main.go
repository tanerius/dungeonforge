package main

import (
	"net"

	"github.com/tanerius/dungeonforge/pkg/logging"
	"github.com/tanerius/dungeonforge/src/lobby"
	"google.golang.org/grpc"
)

func main() {
	service := &lobby.LobbyServerNode{}
	log := logging.NewLogger()

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
