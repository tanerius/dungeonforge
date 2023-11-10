package main

import "github.com/tanerius/dungeonforge/pkg/server"

func main() {
	var server *server.SocketServer = server.NewSocketServer()

	server.StartHTTPServer()
	select {}
}
