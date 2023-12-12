package game

import (
	"time"

	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/database"
	"github.com/tanerius/dungeonforge/pkg/server"

	log "github.com/sirupsen/logrus"
)

type GameServer struct {
	messageChannel chan *GameMessageEvent
	hub            *server.Coordinator
	db             *database.DBWrapper
	em             *eventgoround.EventManager
	players        map[string]chan *GameMessageEvent
}

func NewGameServer(_db *database.DBWrapper, _h *server.Coordinator, _em *eventgoround.EventManager) *GameServer {
	return &GameServer{
		messageChannel: make(chan *GameMessageEvent, 500),
		hub:            _h,
		db:             _db,
		em:             _em,
	}
}

// game server will implement a message handler cuz all it does is process game messages
func (g *GameServer) Type() int {
	return GameMsg
}

func (g *GameServer) HandleEvent(_event *eventgoround.Event) {
	log.Debugln("[GameServer] handling message event")
	if msgEvent, err := eventgoround.GetEventData[*GameMessageEvent](_event); err != nil {
		log.Error(err)
	} else {
		select {
		case g.messageChannel <- msgEvent:
			return
		default:
			log.Debugln("[GameServer] dropped message event")
		}
	}

}

func (g *GameServer) Run() {
	g.em.RegisterListener(g)
	log.Debugln("[GameServer] Running...")
	ticker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case message, ok := <-g.messageChannel:
			if !ok {
				log.Debugln("Main game message channel closed")
				return
			}
			playerChan, playerExists := g.players[message.User.ID.Hex()]
			if !playerExists {
				g.players[message.User.ID.Hex()] = make(chan *GameMessageEvent, 5)
				playerChan = g.players[message.User.ID.Hex()]
				newPlayer := LoadPlayer(message.User, g.db, g.hub)
				go newPlayer.Play(playerChan)
			}
			playerChan <- message
		case <-ticker.C:
			// every 50 ms do samping
		}

	}
}
