package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Identifiable interface {
	GetId() string
	GetHumanId() string
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
