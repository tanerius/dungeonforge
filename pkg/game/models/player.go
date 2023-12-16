package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserId         primitive.ObjectID `bson:"userid,omitempty" json:"-"`
	Gems           int                `bson:"gems,omitempty" json:"gems,omitempty"`
	TotalPurchases int                `bson:"totalpurchases,omitempty" json:"totalpurchases,omitempty"`
	IsNewPlayer    bool               `bson:"isnew,omitempty" json:"isnew,omitempty"`
	Characters     []*Character       `bson:"-" json:"-"`
	IsDirty        bool               `bson:"-" json:"-"`
}
