package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	GameDB          string = "gameDB"
	UsersCollection string = "users"

	// Responses

	RespOK              int = 0
	RespUnknownError    int = 1
	RespConnectionError int = 2
)

// a container for users
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Email        string             `bson:"email,omitempty" json:"-"`
	Password     string             `bson:"password,omitempty" json:"-"`
	Token        string             `bson:"token,omitempty" json:"token"`
	Created      time.Time          `bson:"created,omitempty" json:"created"`
	LastSeen     time.Time          `bson:"lastseen,omitempty" json:"lastseen"`
	IsOnline     bool               `bson:"online,omitempty" json:"-"`
	ClientId     string             `json:"-"`
	ResponseCode int                `json:"responsecode"`
	ResponseMsg  string             `json:"responsemsg"`
}
