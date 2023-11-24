package main

import (
	"log"
	"net/http"

	"github.com/tanerius/dungeonforge/pkg/game"
	"github.com/tanerius/dungeonforge/pkg/server"
)

func main() {
	var gameServer *game.DungeonForge = game.NewDungeonForge()

	var server *server.Server = server.NewServer(gameServer)
	go server.StartServer(true)

	log.Fatal(http.ListenAndServe(":40000", nil))
}
