package entities

import "time"

const (
	GameDB          string = "gameDB"
	UsersCollection string = "users"
)

// a container for users
type User struct {
	*Entity  `bson:",inline"`
	Email    string    `bson:"email,omitempty" json:"email"`
	Password string    `bson:"password,omitempty" json:"password"`
	Token    string    `bson:"token,omitempty" json:"token"`
	Created  time.Time `bson:"Created,omitempty" json:"created"`
	LastSeen time.Time `bson:"Created,omitempty" json:"lastSeen"`
	IsOnline bool      `bson:"online,omitempty" json:"online"`
	ClientId string
}
