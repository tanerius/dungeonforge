package models

type Character struct {
	Name      string `bson:"name" json:"name"`
	Level     int    `bson:"level" json:"level"`
	BaseStats *Stats `bson:"basestats" json:"basestats"`
	ExpStats  *Stats `bson:"expstats" json:"expstats"`
	Gold      int    `bson:"gold" json:"gold"`
	GuildId   string `bson:"guildid" json:"guildid"`
	Exp       int    `bson:"exp" json:"exp"`
}
