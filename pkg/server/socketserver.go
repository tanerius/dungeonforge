package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tanerius/EventGoRound/eventgoround"

	"github.com/gorilla/websocket"
)

type SocketServer struct {
	upgrader     *websocket.Upgrader
	eventManager *eventgoround.EventManager
}

func NewSocketServer(_eventManager *eventgoround.EventManager) *SocketServer {

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
		log.Debugln("[Server] Create connect event ...")
		event := eventgoround.NewEvent(EventClientConnect, NewClientEvent(playerConnection.clientId, playerConnection))
		log.Debugln("[Server] Dispatch connect event ...")
		s.eventManager.DispatchPriorityEvent(event)
	}

	log.Infoln("[Server] Starting HTTP server on port 40000 ...")
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/health", healthCheckHandler)
}
