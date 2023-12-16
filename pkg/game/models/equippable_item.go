package models

type EquippableItem struct {
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
	ItemStats *Modifier `bson:"itemstats" json:"itemstats"`

	Runes []*RuneItem `bson:"runes" json:"runes"`
}

// Implement Id required in the item interface
func (i *EquippableItem) Id() string {
	return i.HrId
}

// Implement Modifiers required in the item interface
func (i *EquippableItem) Modifiers() *Modifier {
	return i.ItemStats
}

// Implement Equip required in the Equippable interface
func (i *EquippableItem) Equip() {
	// TODO: Implement
}

// Implement AddRune required in the Equippable interface
func (i *EquippableItem) AddRune(rune *RuneItem) {
	// TODO: Implement
}
