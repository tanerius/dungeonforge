package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entity
type Entity struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
}

// Create a new entity
func NewEntity(_oid primitive.ObjectID) *Entity {
	return &Entity{
		Id: _oid,
	}
}

func (e *Entity) ID() string {
	if e.Id == primitive.NilObjectID {
		return ""
	}
	return e.Id.Hex()
}
