package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	GameDB          string = "gameDB"
	UsersCollection string = "users"

	// Responses

	RespOK               int = 0
	RespUnknownError     int = 1
	RespConnectionError  int = 2
	RespRegisterError    int = 3
	RespLoginError       int = 4
	RespNotAuthenticated int = 5
	RespInvalidRequest   int = 6
	RespInvalidEmail     int = 7
	RespInvalidUser      int = 8
)

// a container for users
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`
	Email        string             `bson:"email,omitempty" json:"-"`
	Password     string             `bson:"password,omitempty" json:"-"`
	Token        string             `bson:"token,omitempty" json:"token,omitempty"`
	Created      time.Time          `bson:"created,omitempty" json:"created,omitempty"`
	LastSeen     time.Time          `bson:"lastseen,omitempty" json:"lastseen,omitempty"`
	IsOnline     bool               `bson:"online,omitempty" json:"-"`
	ClientId     string             `bson:"-" json:"-"`
	ResponseCode int                `bson:"-" json:"responsecode"`
	ResponseMsg  string             `bson:"-" json:"responsemsg"`
}
