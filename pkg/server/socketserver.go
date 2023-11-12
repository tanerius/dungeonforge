package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type SocketServer struct {
	coord    *Coordinator
	upgrader websocket.Upgrader
}

func NewSocketServer(_c *Coordinator) *SocketServer {
	return &SocketServer{
		coord: _c,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (s *SocketServer) StartHTTPServer() {
	go func() {
		log.Println("SocketServer * Starting HTTP server on port 40000 ...")
		http.HandleFunc("/ws", s.handleWebSocket)
		http.ListenAndServe(":40000", nil)
	}()
}

func (s *SocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("SocketServer * New client trying to connect ...")
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade * ", err)
		return
	}

	playerConnection := newConnection(c)
	log.Println("SocketServer * Registering new client ...")
	s.coord.register <- playerConnection
}
