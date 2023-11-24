package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type Server struct {
	gameServer GameServer
	upgrader   *websocket.Upgrader
}

func NewServer(_gameServer GameServer) *Server {

	return &Server{
		gameServer: _gameServer,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *Server) StartServer(dbg bool) {
	if dbg {
		log.SetLevel(log.DebugLevel)
	}

	healthCheckHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	handleWebSocket := func(w http.ResponseWriter, r *http.Request) {
		log.Debugln("SocketServer * New client trying to connect ...")
		c, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorf("SocketServer upgrade * %v \n", err)
			return
		}

		playerConnection := newClient(c, s)
		log.Debugln("SocketServer * Registering new client ...")

		if err := s.gameServer.HandleClient(playerConnection); err != nil {
			log.Errorf("SocketServer * gameserver said: %v ", err)
			cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error())
			playerConnection.cn.WriteMessage(websocket.CloseMessage, cm)
			playerConnection.cn.Close()
		}

	}

	log.Infoln("SocketServer * Starting HTTP server on port 40000 ...")
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/health", healthCheckHandler)
}
