package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type SocketServer struct {
	gameServer GameServer
	upgrader   websocket.Upgrader
}

func NewSocketServer(_gameServer GameServer) *SocketServer {

	return &SocketServer{
		gameServer: _gameServer,
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
		log.Errorf("SocketServer upgrade * %v \n", err)
		return
	}

	playerConnection := newClient(c)
	log.Println("SocketServer * Registering new client ...")

	if err := s.gameServer.HandleClient(playerConnection); err != nil {
		log.Errorf("SocketServer * gameserver said: %v ", err)
		cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error())
		playerConnection.cn.WriteMessage(websocket.CloseMessage, cm)
		playerConnection.cn.Close()
		close(playerConnection.toSend)
	}

}
