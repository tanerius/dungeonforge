package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/events"

	"github.com/gorilla/websocket"
)

type SocketServer struct {
	upgrader     *websocket.Upgrader
	eventManager *events.EventManager
}

func NewSocketServer(_eventManager *events.EventManager) *SocketServer {

	return &SocketServer{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		eventManager: _eventManager,
	}
}

func (s *SocketServer) StartServer(dbg bool) {
	if dbg {
		log.SetLevel(log.DebugLevel)
	}

	healthCheckHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	handleWebSocket := func(w http.ResponseWriter, r *http.Request) {
		log.Debugln("[Server] New client trying to connect ...")
		c, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorf("SocketServer upgrade * %v \n", err)
			return
		}

		playerConnection := newClient(c, s.eventManager)
		connectionEvent := NewClientEvent(events.EventClientConnected, playerConnection.clientId, playerConnection)
		s.eventManager.Dispatch(connectionEvent)
	}

	log.Infoln("[Server] Starting HTTP server on port 40000 ...")
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/health", healthCheckHandler)
}
