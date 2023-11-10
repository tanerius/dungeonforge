package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/messages"
	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
)

type SocketServer struct {
	server   *gameServer
	upgrader websocket.Upgrader
}

func NewSocketServer() *SocketServer {
	return &SocketServer{
		server: &gameServer{},
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

	defer c.Close()

	log.Println("New client trying to connect")

	_, data, err := c.ReadMessage()

	if err != nil {
		log.Error("read message * ", err)
		return
	}

	var m *messages.Person = &messages.Person{}

	err = proto.Unmarshal(data, m)

	if err != nil {
		log.Error("unmarshal * ", err)
		return
	}

	log.Printf("Name: %s *** Age: %d", m.GetName(), m.GetAge())
}
