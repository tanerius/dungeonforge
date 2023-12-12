package database

type GameDBWrapper struct {
	db *MongoDBWrapper
}

func NewGameDatabase() *GameDBWrapper {
	return &GameDBWrapper{
		db: NewMongoDBWrapper("mongodb://dungeonmaster:m123123123@localhost:27017/", 100),
	}
}
