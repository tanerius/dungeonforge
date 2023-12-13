package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// an interface used for any item
type Identifiable interface {
	GetId() string
	GetHumanId() string
}

// an interface for retailableitems
type Sellable interface {
	Sell(*Character) (int, error)
}

// an interface used for potions
type Consumable interface {
	Consume(*Character) (time.Duration, error)
	TimeRemaining() time.Duration
	IsConsumed() bool
}

// an interface used for potions
type Equipable interface {
	Equip(*Character) (string, error)
	UnEquip(*Character) error
}

type Entity struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	// A human readable unique ID which will be used for localization later. Index also
	HrId string `bson:"hrid" json:"hrid"`

	// Name of an item in english
	Name string `bson:"name" json:"name"`

	// Item description in english
	Desc string `bson:"desc" json:"desc"`
}
