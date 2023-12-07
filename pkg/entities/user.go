package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	GameDB          string = "gameDB"
	UsersCollection string = "users"
)

// a container for users
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Email    string             `bson:"email,omitempty" json:"email"`
	Password string             `bson:"password,omitempty" json:"password"`
	Token    string             `bson:"token,omitempty" json:"token"`
	Created  time.Time          `bson:"created,omitempty" json:"created"`
	LastSeen time.Time          `bson:"lastseen,omitempty" json:"lastseen"`
	IsOnline bool               `bson:"online,omitempty" json:"online"`
	ClientId string             `json:"-"`
}
