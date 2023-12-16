package models

import (
	"github.com/tanerius/dungeonforge/pkg/game/gameobjects"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Character struct {
	Id            primitive.ObjectID             `bson:"_id,omitempty" json:"-"`
	PlayerId      primitive.ObjectID             `bson:"playerid,omitempty" json:"-"`
	Name          string                         `bson:"name" json:"name"`
	Level         int                            `bson:"level" json:"level"`
	BaseStats     *Stats                         `bson:"basestats" json:"basestats"`
	ExpStats      *Stats                         `bson:"expstats" json:"expstats"`
	Gold          int                            `bson:"gold" json:"gold"`
	GuildId       string                         `bson:"guildid" json:"guildid"`
	Exp           int                            `bson:"exp" json:"exp"`
	Inventory     map[string]Item                `bson:"inventory" json:"inventory"`
	EquippedItems [gameobjects.SlotMaxVal]string `bson:"equipped" json:"equipped"`
}
