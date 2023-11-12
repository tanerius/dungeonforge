package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/messages"

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

	// Handle incoming messages from the player.
	go func() {
		defer func() {
			log.Printf("SocketServer * Deregistering client %s \n", playerConnection.entityId.String())
			s.coord.unregister <- playerConnection
		}()

		for {
			// Read messages from the player and handle them as needed.
			// _, message, err := playerConnection.cn.ReadMessage() // Use for protobuff
			var message messages.Payload
			err := playerConnection.cn.ReadJSON(&message)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseMessage) {
					log.Printf("Client closed connection: %v ", err)
				} else {
					log.Printf("error: %v", err)
				}
				return
			}

			log.Printf("SocketServer * Got message from client %s \n", playerConnection.entityId.String())

			// Send the received message to the game loop for processing.
			s.coord.playerMessages <- &message
		}
	}()
}
