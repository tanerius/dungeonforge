package game

import (
	"time"

	"github.com/tanerius/dungeonforge/pkg/entities"
	"github.com/tanerius/dungeonforge/pkg/game/models"
	"github.com/tanerius/dungeonforge/pkg/server"

	log "github.com/sirupsen/logrus"
)

type GameInstance struct {
	user   *entities.User
	db     *GameDBWrapper
	hub    *server.Coordinator
	player *models.Player
}

func SpawnInstance(_u *entities.User, _db *GameDBWrapper, _hub *server.Coordinator) *GameInstance {
	return &GameInstance{
		user: _u,
		db:   _db,
		hub:  _hub,
	}
}

func (g *GameInstance) Play(msgChan <-chan *GameMessageEvent) {
	log.Debugf("Player:Playing %s", g.user.Name)

	player, err := g.db.GetPlayer(g.user.ID)
	if err != nil {
		log.Error(err)
		return
	}

	g.player = player

	characters, err := g.db.GetPlayerCharacters(player.Id)
	if err != nil {
		log.Error(err)
		return
	}

	g.player.Characters = characters

	defer func() {
		log.Debugf("Player:Stopping %s", g.user.Name)
	}()

	ticker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return
			}
			log.Debugf("Player: %s provessing a message", g.player)
			log.Debugf("%v", msg)
		case <-ticker.C:
			// every 100 ms do samping
		}
	}
}
