package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CharacterRace type
type CharacterRace int

const (
	crHuman CharacterRace = iota + 1
	crDwarf
	crElve
	crElemental
	crMonster
	crFantasy
)

func (cr CharacterRace) String() string {
	return [...]string{"human", "dwarf", "elve", "elemental", "crMonster"}[cr]
}

//Party data
type Party struct {
	*Entity    `bson:",inline"`
	Name       string               `json:"name"`
	Created    time.Time            `bson:"created,omitempty" json:"created,omitempty"`
	Characters []primitive.ObjectID `bson:"characters,omitempty" json:"characters,omitempty"`
}
