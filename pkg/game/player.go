package game

import (
	"time"

	"github.com/tanerius/dungeonforge/pkg/database"
	"github.com/tanerius/dungeonforge/pkg/entities"
	"github.com/tanerius/dungeonforge/pkg/server"

	log "github.com/sirupsen/logrus"
)

type Player struct {
	user *entities.User
	db   *database.DBWrapper
	hub  *server.Coordinator
}

func LoadPlayer(_u *entities.User, _db *database.DBWrapper, _hub *server.Coordinator) *Player {
	return &Player{
		user: _u,
		db:   _db,
		hub:  _hub,
	}
}

func (p *Player) Play(msgChan <-chan *GameMessageEvent) {
	log.Debugf("Player:Playing %s", p.user.Name)

	defer func() {
		log.Debugf("Player:Stopping %s", p.user.Name)
	}()

	ticker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return
			}
			log.Debugf("%v", msg)
		case <-ticker.C:
			// every 100 ms do samping
		}
	}
}
