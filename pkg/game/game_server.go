package game

import (
	"encoding/json"
	"os"
	"time"

	"github.com/tanerius/EventGoRound/eventgoround"
	"github.com/tanerius/dungeonforge/pkg/game/models"
	"github.com/tanerius/dungeonforge/pkg/server"

	log "github.com/sirupsen/logrus"
)

type GameServer struct {
	playerLoopStop chan string
	messageChannel chan *GameMessageEvent
	hub            *server.Coordinator
	db             *GameDBWrapper
	em             *eventgoround.EventManager
	players        map[string]chan *GameMessageEvent
	items          []models.Equippable
	potions        []models.Consumable
}

func NewGameServer(_h *server.Coordinator, _em *eventgoround.EventManager) *GameServer {
	return &GameServer{
		messageChannel: make(chan *GameMessageEvent, 500),
		hub:            _h,
		db:             NewGameDatabase(),
		em:             _em,
		players:        make(map[string]chan *GameMessageEvent, 100),
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
	g.items = make([]models.Equippable, 0)
	g.potions = make([]models.Consumable, 0)

	g.readItems()
	g.em.RegisterListener(g)
	log.Debugln("[GameServer] Running...")
	ticker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case pid, ok := <-g.playerLoopStop:
			if !ok {
				log.Debugln("Loop stop channel closed")
				return
			}

			if playerChan, playerExists := g.players[pid]; playerExists {
				close(playerChan)
				delete(g.players, pid)
			}

		case message, ok := <-g.messageChannel:
			if !ok {
				log.Debugln("Main game message channel closed")
				return
			}
			playerChan, playerExists := g.players[message.User.ID.Hex()]
			if !playerExists {
				log.Debugf("\n%v\n", message)
				g.players[message.User.ID.Hex()] = make(chan *GameMessageEvent, 5)
				playerChan = g.players[message.User.ID.Hex()]
				newPlayer := SpawnInstance(g, message.User, g.db, g.hub)
				go newPlayer.Play(playerChan)
			}
			playerChan <- message
		case <-ticker.C:
			// every 50 ms do samping
		}

	}
}

// this function reads items one by one
func (g *GameServer) readItems() {

	// weapons
	weapons, err := os.Open("pkg/game/fixtures/main_weapons.json")
	if err != nil {
		log.Panic("opening json file", err.Error())
	}

	jsonParser := json.NewDecoder(weapons)
	if err = jsonParser.Decode(&g.items); err != nil {
		log.Panic("parsing config file", err.Error())
	}

	weapons.Close()

	// potions
	potions, err := os.Open("pkg/game/fixtures/potions.json")
	if err != nil {
		log.Panic("opening json file", err.Error())
	}

	jsonParser = json.NewDecoder(potions)

	p := make([]models.Consumable, 0)

	if err = jsonParser.Decode(&p); err != nil {
		log.Panic("parsing config file", err.Error())
	}

	g.potions = append(g.potions, p...)
	potions.Close()

}
