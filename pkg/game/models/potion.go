package models

type Potion struct {
	Entity     `bson:",inline" json:",inline"`
	MarketInfo `bson:",inline" json:",inline"`

	// Indicates if an item can be dropped in game
	Droppable bool `bson:"droppable" json:"droppable"`

	// How the item changes the wielder's stats
	StatsWielder *Stats `bson:"statswielder" json:"statswielder"`

	// Duration of the potion's effects in seconds.
	//
	// 0 means once active it remains active
	Duration uint `bson:"duration" json:"duration"`
}
