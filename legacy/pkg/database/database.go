package database

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tanerius/dungeonforge/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (d *DBWrapper) Logout(_hexid string) error {
	objectId, errHex := primitive.ObjectIDFromHex(_hexid)
	if errHex != nil {
		return errHex
	}
	update := bson.M{
		"$set": bson.M{"token": "", "online": false},
	}

	_, err := d.db.UpdateDocument(entities.GameDB, entities.UsersCollection,
		bson.M{"_id": objectId}, update)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBWrapper) Login(name, pass string) (*entities.User, error) {
	// 6) Create the update
	newToken := uuid.NewString()
	now := time.Now()
	update := bson.M{
		"$set": bson.M{"lastSeen": now, "token": newToken, "online": true},
	}

	result, err := d.db.GetDocumentWithUpdate(entities.GameDB, entities.UsersCollection, bson.M{"name": name, "password": GetMD5Hash(pass)}, update)
	if err != nil {
		return nil, err
	}

	var retrievedDoc *entities.User = &entities.User{}

	if err := result.Decode(retrievedDoc); err != nil {
		return nil, err
	}
	retrievedDoc.Token = newToken
	retrievedDoc.LastSeen = now

	return retrievedDoc, nil
}

func (d *DBWrapper) Register(email, pass, name string) (*entities.User, error) {
	hashedPass := GetMD5Hash(pass)

	filter := bson.M{
		"$or": bson.A{
			bson.M{"email": email},
			bson.M{"name": name},
		},
	}

	count := d.db.Exists(entities.GameDB, entities.UsersCollection, filter)
	if count > 0 {
		return nil, errors.New("user already exists")
	}

	newUser := &entities.User{
		Email:    email,
		Name:     name,
		Password: hashedPass,
		Token:    uuid.NewString(),
		Created:  time.Now(),
		LastSeen: time.Now(),
		IsOnline: true,
	}

	inserted, err := d.db.CreateDocument(entities.GameDB, entities.UsersCollection, newUser)
	if err != nil {
		return nil, err
	} else {
		if oid, ok := inserted.InsertedID.(primitive.ObjectID); ok {
			newUser.ID = oid
		}
		return newUser, nil
	}
}
