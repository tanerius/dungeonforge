package game

import "github.com/tanerius/dungeonforge/pkg/database"

type GameDBWrapper struct {
	db *database.MongoDBWrapper
}

func NewGameDatabase() *GameDBWrapper {
	return &GameDBWrapper{
		db: database.NewMongoDBWrapper("mongodb://dungeonmaster:m123123123@localhost:27017/", 100),
	}
}
