package models

type Item struct {
	Entity     `bson:",inline" json:",inline"`
	MarketInfo `bson:",inline" json:",inline"`

	// Indicates if an item can be dropped in game
	Droppable bool `bson:"droppable" json:"droppable"`

	// can an item be buffed by a rune
	Buffable bool `bson:"buffable" json:"buffable"`

	// Which slot the item occupies when equipped
	Slot int `bson:"slot" json:"slot"`

	// Rarity of an item
	Rarity int `bson:"rarity" json:"rarity"`

	// Minimum level required to equip item
	MinLevel int `bson:"minlevel" json:"minlevel"`

	// On which level should we stop dropping this item.
	StopDropLevel int `bson:"stopdrop" json:"-"`

	// How the item changes the wielder's stats
	StatsWielder *Stats `bson:"statswielder" json:"statswielder"`

	// How the item changes the enemy's stats when equipped by the wielder.
	// These are for items that not only buff the wielder but also debuff the enemy.
	//
	// Example: To reduce 10% of enemy strength, set StatsEnemy.Str = -0.1
	// Example: To increase enemy's dexterity by 60, set StatsEnemy.Str = 60.0
	StatsEnemy *Stats `bson:"statsenemy" json:"statsenemy"`

	// Damage modifier of the item
	*Damage `bson:"damage" json:"damage"`
}
