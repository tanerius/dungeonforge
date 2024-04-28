package game

import (
	"time"

	"github.com/tanerius/dungeonforge/pkg/entities"
	"github.com/tanerius/dungeonforge/pkg/server"

	log "github.com/sirupsen/logrus"
)

type GameInstance struct {
	user *entities.User
	db   *GameDBWrapper
	hub  *server.Coordinator
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
			log.Debugf("Player: %s provessing a message", g.user.Name)
			log.Debugf("%v", msg)
		case <-ticker.C:
			// every 100 ms do samping
		}
	}
}
