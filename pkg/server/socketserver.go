package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type SocketServer struct {
	server   *GameServer
	coord    *Coordinator
	upgrader websocket.Upgrader
}

func NewSocketServer(_gs *GameServer, _c *Coordinator) *SocketServer {
	return &SocketServer{
		server: _gs,
		coord:  _c,
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
	log.Println("New client trying to connect...")
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade * ", err)
		return
	}

	playerConnection := newConnection(c)
	s.coord.register <- playerConnection

	// Handle incoming messages from the player.
	go func() {
		defer func() {
			s.coord.unregister <- playerConnection
			playerConnection.cn.Close()
		}()

		for {
			// Read messages from the player and handle them as needed.
			_, message, err := playerConnection.cn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			// Process player input, update game state, and send responses.
			// Example: player.conn.WriteMessage(websocket.TextMessage, responseBytes)

			// Send the received message to the game loop for processing.
			s.coord.playerMessages <- message
		}
	}()
}
