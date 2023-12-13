package game

import (
	"github.com/tanerius/dungeonforge/pkg/database"
	"github.com/tanerius/dungeonforge/pkg/game/gameobjects"
	"github.com/tanerius/dungeonforge/pkg/game/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameDBWrapper struct {
	db *database.MongoDBWrapper
}

func NewGameDatabase() *GameDBWrapper {
	return &GameDBWrapper{
		db: database.NewMongoDBWrapper(gameobjects.GameDbURI, gameobjects.ConnPoolMaxSize),
	}
}

func (d *GameDBWrapper) GetPlayer(_id primitive.ObjectID) (*models.Player, error) {
	result, err := d.db.GetDocument(gameobjects.GameDB, gameobjects.ColPlayer, bson.M{"_id": _id})
	if err != nil {
		return nil, err
	}

	var retrievedDoc *models.Player = &models.Player{}

	if err := result.Decode(retrievedDoc); err != nil {
		return nil, err
	}

	return retrievedDoc, nil
}
