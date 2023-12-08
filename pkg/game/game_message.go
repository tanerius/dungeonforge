package game

import "github.com/tanerius/dungeonforge/pkg/entities"

// Definition of game messages
type GameMessageEvent struct {
	ClientId string
	User     *entities.User
	Data     map[string]string
}
