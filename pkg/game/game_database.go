package game

import (
	"context"
	"errors"

	"github.com/tanerius/dungeonforge/pkg/database"
	"github.com/tanerius/dungeonforge/pkg/entities"
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
		// player with id doesnt exist creape a new one

		filter := bson.M{"userid": _id}
		count := d.db.Exists(gameobjects.GameDB, gameobjects.ColPlayer, filter)
		if count > 0 {
			return nil, errors.New("player with given userID already exists")
		}

		newPlayer := &models.Player{
			UserId:         _id,
			Gems:           gameobjects.StartingGems,
			TotalPurchases: 0,
			IsNewPlayer:    true,
			Characters:     make([]*models.Character, 0),
			IsDirty:        false,
		}

		inserted, err := d.db.CreateDocument(entities.GameDB, entities.UsersCollection, newPlayer)
		if err != nil {
			return nil, err
		} else {
			if oid, ok := inserted.InsertedID.(primitive.ObjectID); ok {
				newPlayer.Id = oid
			}
			return newPlayer, nil
		}
	}

	var retrievedDoc *models.Player = &models.Player{}

	if err := result.Decode(retrievedDoc); err != nil {
		return nil, err
	}

	return retrievedDoc, nil
}

func (d *GameDBWrapper) SavePlayer(_player *models.Player) (*models.Player, error) {

	_, err := d.db.UpdateDocument(gameobjects.GameDB, gameobjects.ColPlayer, bson.M{"_id": _player.Id}, _player)
	if err != nil {
		return _player, err
	}

	_player.IsDirty = false

	return _player, nil
}

func (d *GameDBWrapper) GetPlayerCharacters(_pid primitive.ObjectID) ([]*models.Character, error) {
	ret := make([]*models.Character, 0)
	cursor, err := d.db.GetDocuments(gameobjects.GameDB, gameobjects.ColCharacter, bson.M{"playerid": _pid})
	if err != nil {
		return ret, err
	}

	if err := cursor.All(context.TODO(), &ret); err != nil {
		return ret, err
	}

	return ret, nil
}
