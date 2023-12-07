package database

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/tanerius/dungeonforge/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

type DBWrapper struct {
	db *MongoDBWrapper
}

func NewDatabase() *DBWrapper {
	return &DBWrapper{
		db: NewMongoDBWrapper("mongodb://dungeonmaster:m123123123@localhost:27017/", 100),
	}
}

func (d *DBWrapper) Login(email, pass string) (*entities.User, error) {
	result, err := d.db.GetDocument(entities.GameDB, entities.UsersCollection, bson.M{"email": email, "pass": GetMD5Hash(pass)})
	if err != nil {
		return nil, err
	}

	var retrievedDoc *entities.User = &entities.User{}

	if err := result.Decode(retrievedDoc); err != nil {
		return nil, err
	}

	return retrievedDoc, nil
}

func (db *DBWrapper) Register(email, pass string) string {
	return ""
}
