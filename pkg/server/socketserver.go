package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type SocketServer struct {
	upgrader websocket.Upgrader
}

func NewSocketServer() *SocketServer {
	return &SocketServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (s *SocketServer) StartHTTPServer() {
	go func() {
		log.Println("Starting HTTP server on port 40000 ...")
		http.HandleFunc("/ws", s.handleWebSocket)
		http.ListenAndServe(":40000", nil)
	}()
}

func (s *SocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade * ", err)
		return
	}

	log.Println("New client trying to connect")
	log.Print(c)
}
